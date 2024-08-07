package mesos

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	mesosaws "github.com/AVENTER-UG/mesos-autoscale/aws"
	"github.com/AVENTER-UG/mesos-autoscale/redis"
	cfg "github.com/AVENTER-UG/mesos-autoscale/types"
	util "github.com/AVENTER-UG/util/util"
	"github.com/sirupsen/logrus"
)

// Scheduler include all the current vars and global config
type Scheduler struct {
	Config *cfg.Config
	Redis  *redis.Redis
	AWS    *mesosaws.AWS
	Health bool
}

// New will create the Scheduler object
func New(cfg *cfg.Config) *Scheduler {
	e := &Scheduler{
		Config: cfg,
		Health: true,
	}
	// Add protocoll to the endpoint depends if SSL is enabled
	logrus.Info("Connect Provider Apache Mesos Airflow Provider: ", e.Config.AirflowMesosScheduler)

	return e
}

// EventLoop is the main loop of all events
func (e *Scheduler) EventLoop() {
	ticker := time.NewTicker(e.Config.PollInterval)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		if e.Health {
			// connect to the scheduler and read all dags
			e.saveDags()

			// remove dags from redis if it's out of the queue
			e.delDags()

			// check if we have to scale out instances
			e.checkDags()

			// Check if we can terminate the ec2 instances
			go e.checkEC2Instance()
		}
	}
}

// delDags will delete the dags from redis if it's not in the queue anymore
func (e *Scheduler) delDags() {
	keys := e.Redis.GetAllRedisKeys(e.Config.RedisPrefix + ":dags:*")
	for keys.Next(e.Redis.RedisCTX) {
		i := e.Redis.GetTaskFromRunID(keys.Val())
		if i != nil {
			if !e.dagIsIn(i) {
				ret := e.Redis.DelDagTaskRedis(*i)
				logrus.WithField("func", "scheduler.delDags").Debugf("Remove DAG (%s, %s) from redis: %d ", i.TaskID, i.RunID, ret)
			}
		}
	}
}

// check if the give das is still in queue
func (e *Scheduler) dagIsIn(dag *cfg.DagTask) bool {
	dags := e.getDags()
	if dags == nil {
		return false
	}
	for _, i := range dags {
		if dag.RunID == i.RunID {
			return true
		}
	}

	return false
}

// HealthCheck will check the health
func (e *Scheduler) HealthCheck() {
	ticker := time.NewTicker(e.Config.PollInterval)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		err := e.Redis.PingRedis()
		if err != nil {
			logrus.WithField("func", "mesos.HealthCheck").Error("Redis connection error:", err.Error())
			e.Health = false
		} else {
			e.Health = true
		}
	}
}

// check if we have to scale out instances
func (e *Scheduler) checkDags() {
	keys := e.Redis.GetAllRedisKeys(e.Config.RedisPrefix + ":dags:*")
	for keys.Next(e.Redis.RedisCTX) {
		i := e.Redis.GetTaskFromRunID(keys.Val())
		if i != nil {
			timeDiff := time.Since(i.StartDate).Seconds()
			if timeDiff >= e.Config.WaitTimeout.Seconds() && !i.ASG {
				e.scaleOut(i, i.MesosExecutor.InstanceType)
			}
		}
	}
}

func (e *Scheduler) scaleOut(task *cfg.DagTask, instanceType string) {
	logrus.WithField("func", "checkDags").Info("ScaleOut Mesos: ", task.TaskID)
	task.ASG = true
	e.Redis.SaveDagTaskRedis(*task)

	var ec cfg.EC2Struct
	if instanceType == "" {
		mem := int64(e.convertMemoryToFloat(task.MesosExecutor.MemLimit) * 1.1)
		cpu := int64(task.MesosExecutor.Cpus)

		logrus.WithField("func", "mesos.scaleOut").Tracef("Need Mem: %d ", mem)
		logrus.WithField("func", "mesos.scaleOut").Tracef("Need CPU: %d ", cpu)
		logrus.WithField("func", "mesos.scaleOut").Tracef("Need Architecture: %s ", task.MesosExecutor.Architecture)

		instanceType = e.AWS.FindMatchedInstanceType(mem, cpu, task.MesosExecutor.Architecture)
	}
	ec.EC2 = e.AWS.CreateInstance(instanceType)
	e.Redis.SaveEC2InstanceRedis(ec)
}

// get all running dags from airflow mesos scheduler
func (e *Scheduler) saveDags() {
	dags := e.getDags()

	for _, i := range dags {
		sTask := e.Redis.GetTaskFromRunID(e.Config.RedisPrefix + ":dags:" + i.DagID + ":" + i.TaskID + ":" + i.RunID + ":" + strconv.Itoa(i.TryNumber))
		if sTask == nil {
			mem := e.convertMemoryToFloat(i.MesosExecutor.MemLimit)
			if i.MesosExecutor.InstanceType == "" {
				i.MesosExecutor.InstanceType = e.AWS.FindMatchedInstanceType(int64(mem*1.1), int64(i.MesosExecutor.Cpus), i.MesosExecutor.Architecture)
			}
			i.StartDate = time.Now()
			e.Redis.SaveDagTaskRedis(i)
			logrus.WithField("func", "EventLoop").Info("Found new DAG in queue: ", i.DagID)
			logrus.WithField("func", "EventLoop").Trace("Dag ID: ", i.DagID)
			logrus.WithField("func", "EventLoop").Trace("Dag Task ID: ", i.TaskID)
			logrus.WithField("func", "EventLoop").Trace("Dag Run ID: ", i.RunID)
			logrus.WithField("func", "EventLoop").Trace("Dag TryNumber: ", strconv.Itoa(i.TryNumber))
			logrus.WithField("func", "EventLoop").Trace("Dag StartDate: ", i.StartDate)
			logrus.WithField("func", "EventLoop").Trace("Dag CPUs: ", i.MesosExecutor.Cpus)
			logrus.WithField("func", "EventLoop").Trace("Dag MEM: ", mem)
			logrus.WithField("func", "EventLoop").Trace("Dag Architecture: ", i.MesosExecutor.Architecture)
			logrus.WithField("func", "EventLoop").Trace("Dag InstanceType: ", i.MesosExecutor.InstanceType)
			logrus.WithField("func", "EventLoop").Trace("ASG: ", i.ASG)
			logrus.WithField("func", "EventLoop").Trace("---------------------------------------")
		} else if sTask.ASG {
			// if the task is already in the queue, marked but still not running,
			// then search a matching instance type
			timeDiff := time.Since(sTask.StartDate).Seconds()
			if timeDiff >= e.Config.WaitTimeoutOverwrite.Seconds() {
				logrus.WithField("func", "EventLoop").Debugf("DAG (%s) still not running. Try other instance type: ", i.DagID)
				sTask.StartDate = time.Now()
				e.scaleOut(sTask, "")
			}
		}
	}
}

// get all running dags from airflow mesos scheduler
func (e *Scheduler) getDags() []cfg.DagTask {
	client := &http.Client{}
	client.Transport = &http.Transport{
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("GET", e.Config.AirflowMesosScheduler+"/v0/dags", nil)
	req.Close = true
	req.SetBasicAuth(e.Config.APIUsername, e.Config.APIPassword)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logrus.WithField("func", "getDags").Error("Could not get Dags from airflow: ", err.Error())
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logrus.WithField("func", "getDags").Error("Status Code not 200: ", res.StatusCode)
		return nil
	}

	var dags []cfg.DagTask
	err = json.NewDecoder(res.Body).Decode(&dags)
	if err != nil {
		logrus.WithField("func", "getDags").Error("Cannot decode json: ", err.Error())
		return nil
	}

	return dags
}

func (e *Scheduler) convertMemoryToFloat(memoryStr string) float64 {
	memoryStr = strings.ToLower(memoryStr)
	var memoryVal float64
	var err error

	if strings.HasSuffix(memoryStr, "g") {
		memoryVal, err = strconv.ParseFloat(memoryStr[:len(memoryStr)-1], 64)
		if err != nil {
			return 1000.0
		}
		memoryVal *= 1024
	} else if strings.HasSuffix(memoryStr, "m") {
		memoryVal, err = strconv.ParseFloat(memoryStr[:len(memoryStr)-1], 64)
		if err != nil {
			return 1000.0
		}
	}

	return memoryVal
}

// check if the ec2 instance still running mesos tasks
func (e *Scheduler) checkEC2Instance() {
	keys := e.Redis.GetAllRedisKeys(e.Config.RedisPrefix + ":ec2:*")
	// instance count for summary
	i := 0
	// instances in check mode for summary
	c := 0
	for keys.Next(e.Redis.RedisCTX) {
		i++
		instance := e.Redis.GetEC2InstanceFromID(keys.Val())
		if instance == nil {
			return
		}
		if len(instance.EC2.Instances) <= 0 {
			logrus.WithField("func", "checkEC2Instance").Debug("Didnt got instances")
			continue
		}

		// Only check if there is not already a check running
		if !instance.Check {
			// if the launch time is to short, do not try to connect reach the agent
			// get current time
			timeNow := time.Now()
			timeDiff := timeNow.Sub(*instance.EC2.Instances[0].LaunchTime).Minutes()

			if timeDiff <= e.Config.AWSTerminateWait.Minutes() {
				continue
			}

			c++

			// set check state to true
			instance.Check = true
			e.Redis.SaveEC2InstanceRedis(*instance)

			client := &http.Client{
				Timeout: e.Config.MesosAgentTimeout,
			}
			client.Transport = &http.Transport{
				// #nosec G402
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			protocol := "https"
			if !e.Config.MesosAgentSSL {
				protocol = "http"
			}

			hostIP := instance.EC2.Instances[0].NetworkInterfaces[0].PrivateIpAddress

			req, _ := http.NewRequest("POST", protocol+"://"+*hostIP+":"+e.Config.MesosAgentPort+"/state", nil)
			req.Close = true
			req.SetBasicAuth(e.Config.MesosAgentUsername, e.Config.MesosAgentPassword)
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)

			if err != nil {
				logrus.WithField("func", "checkEC2Instances").Error("Could not connect to agent: ", err.Error())
				instance.Check = false
				e.AWS.TerminateInstance(instance.EC2.Instances[0].InstanceId)
				e.Redis.DelRedisKey(e.Config.RedisPrefix + ":ec2:" + *instance.EC2.Instances[0].InstanceId)
				// create a new instance
				var ec cfg.EC2Struct
				ec.EC2 = e.AWS.CreateInstance(*instance.EC2.Instances[0].InstanceType)
				e.Redis.SaveEC2InstanceRedis(ec)
				continue
			}

			// set mesos agent error counter to zero if we got a connection
			instance.AgentTimeout = 0

			defer res.Body.Close()

			if res.StatusCode != 200 {
				logrus.WithField("func", "checkEC2Instances").Error("HTTP status code: ", res.StatusCode)
				instance.Check = false
				e.Redis.SaveEC2InstanceRedis(*instance)
				continue
			}

			var agent cfg.MesosAgentState
			err = json.NewDecoder(res.Body).Decode(&agent)
			if err != nil {
				logrus.WithField("func", "checkEC2Instances").Error("Could not decode json result: ", err.Error())
				// Den Inhalt des Body als Text ausgeben
				body, err := io.ReadAll(res.Body)
				if err == nil {
					logrus.WithField("func", "checkEC2Instances").Error("Error URL: " + protocol + "://" + *hostIP + ":" + e.Config.MesosAgentPort + "/state")
					logrus.WithField("func", "checkEC2Instances").Error("Error Body: ", string(body))
				}
				// uncheck these instance
				instance.Check = false
				e.Redis.SaveEC2InstanceRedis(*instance)
				continue
			}

			// uncheck these instance
			instance.Check = false
			e.Redis.SaveEC2InstanceRedis(*instance)

			// check if there is still a airflow task running or if the instance reach the max age
			if (!e.isFrameworkName(agent) || timeDiff >= e.Config.AWSInstanceMaxAge.Minutes()) && e.Config.AWSInstanceTerminate {
				e.terminateNode(agent, instance)
				e.deactivateNode(agent)
			}
		}
	}
	if i > 0 {
		logrus.WithField("func", "mesos.checkEC2Instance").Debugf("There are %d instances in DB. %d of them are in check mode.", i, c)
	}
}

func (e *Scheduler) deactivateNode(agent cfg.MesosAgentState) {
	logrus.WithField("func", "mesos.deactivateNode").Info("Deactivate Node: ", agent.ID)

	var deactivateCall cfg.MesosAgentDeactivate
	deactivateCall.Type = "DEACTIVATE_AGENT"
	deactivateCall.DeactivateAgent.AgentID.Value = agent.ID

	d, err := json.Marshal(deactivateCall)
	if err != nil {
		logrus.WithField("func", "mesos.deactivateNode").Error("Could not encode json: ", err.Error())
	}

	protocol := "https"
	if !e.Config.MesosAgentSSL {
		protocol = "http"
	}
	client := &http.Client{
		Timeout: e.Config.MesosAgentTimeout,
	}
	client.Transport = &http.Transport{
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("POST", protocol+"://"+e.Config.MesosMaster+":"+e.Config.MesosMasterPort+"/api/v1", bytes.NewBuffer([]byte(d)))
	req.Close = true
	req.SetBasicAuth(e.Config.MesosMasterUsername, e.Config.MesosMasterPassword)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if res.StatusCode != 200 {
		logrus.WithField("func", "mesos.deactivateNode").Error("Could not deactivate. Status code: ", res.StatusCode)
		logrus.WithField("func", "mesos.deactivateNode").Debug("JSON: ", util.PrettyJSON(d))
	}

	defer res.Body.Close()
	if err != nil {
		logrus.WithField("func", "mesos.deactivateNode").Error("Could not connect to agent: ", err.Error())
	}
}

func (e *Scheduler) terminateNode(agent cfg.MesosAgentState, instance *cfg.EC2Struct) {

	protocol := "https"
	if !e.Config.MesosAgentSSL {
		protocol = "http"
	}
	client := &http.Client{
		Timeout: e.Config.MesosAgentTimeout,
	}
	client.Transport = &http.Transport{
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("POST", protocol+"://"+e.Config.MesosMaster+":"+e.Config.MesosMasterPort+"/slaves/"+agent.ID, nil)
	req.Close = true
	req.SetBasicAuth(e.Config.MesosMasterUsername, e.Config.MesosMasterPassword)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logrus.WithField("func", "terminateNode").Error("Could not connect to agent: ", err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var state cfg.MesosAgent
		err = json.NewDecoder(res.Body).Decode(&state)
		if err != nil {
			logrus.WithField("func", "terminateNode").Error("Could not decode json result: ", err.Error())
			// if there is an error, dump out the res.Body as debug
			bodyBytes, err := io.ReadAll(res.Body)
			if err == nil {
				logrus.WithField("func", "terminateNode").Debug("response Body Dump: ", string(bodyBytes))
			}
		}

		// get the used agent info
		for _, a := range state.Slaves {
			if a.ID == agent.ID && a.Deactivated {
				logrus.WithField("func", "terminateNode").Info("Terminate Node: ", agent.ID)
				e.AWS.TerminateInstance(instance.EC2.Instances[0].InstanceId)
				e.Redis.DelRedisKey(e.Config.RedisPrefix + ":ec2:" + *instance.EC2.Instances[0].InstanceId)
			}
		}
	}
}

// check if there is still a airflow job running
func (e *Scheduler) isFrameworkName(agent cfg.MesosAgentState) bool {
	for _, framework := range agent.Frameworks {
		if strings.EqualFold(framework.Name, e.Config.AirflowMesosName) {
			return true
		}
	}
	return false
}
