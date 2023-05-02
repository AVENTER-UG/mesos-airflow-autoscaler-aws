package mesos

import (
	"crypto/tls"
	cTls "crypto/tls"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	mesosaws "github.com/AVENTER-UG/mesos-autoscale/aws"
	"github.com/AVENTER-UG/mesos-autoscale/redis"
	cfg "github.com/AVENTER-UG/mesos-autoscale/types"
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
			e.getDags()

			// check if we have to scale out instances
			e.checkDags()

			// Check if we can terminate the ec2 instances
			go e.checkEC2Instance()
		}
	}
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

		timeDiff := time.Since(i.StartDate).Seconds()
		if timeDiff >= e.Config.WaitTimeout.Seconds() && !i.ASG {
			logrus.WithField("func", "checkDags").Info("ScaleOut Mesos: ", i.RunID)
			i.ASG = true
			e.Redis.SaveDagTaskRedis(*i)

			var ec cfg.EC2Struct
			if i.MesosExecutor.InstanceType != "" {
				ec.EC2 = e.AWS.CreateInstance(i.MesosExecutor.InstanceType)
				e.Redis.SaveEC2InstanceRedis(ec)
			} else {
				if i.MesosExecutor.MemLimit >= 32768 {
					ec.EC2 = e.AWS.CreateInstance(e.Config.AWSInstance64)
					e.Redis.SaveEC2InstanceRedis(ec)
				} else if i.MesosExecutor.MemLimit >= 16384 {
					ec.EC2 = e.AWS.CreateInstance(e.Config.AWSInstance32)
					e.Redis.SaveEC2InstanceRedis(ec)
				} else {
					ec.EC2 = e.AWS.CreateInstance(e.Config.AWSInstance16)
					e.Redis.SaveEC2InstanceRedis(ec)
				}
			}
		}
	}
}

// get all running dags from airflow mesos scheduler
func (e *Scheduler) getDags() {
	client := &http.Client{}
	client.Transport = &http.Transport{
		// #nosec G402
		TLSClientConfig: &cTls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("GET", e.Config.AirflowMesosScheduler+"/v0/dags", nil)
	req.Close = true
	req.SetBasicAuth(e.Config.APIUsername, e.Config.APIPassword)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logrus.WithField("func", "getDags").Error("Could not get Dags from airflow: ", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logrus.WithField("func", "getDags").Error("Status Code not 200: ", res.StatusCode)
		return
	}

	var dags []cfg.DagTask
	err = json.NewDecoder(res.Body).Decode(&dags)
	if err != nil {
		logrus.WithField("func", "getDags").Error("Cannot decode json: ", err.Error())
		return
	}

	for _, i := range dags {
		sTask := e.Redis.GetTaskFromRunID(e.Config.RedisPrefix + ":dags:" + i.DagID + ":" + i.TaskID + ":" + i.RunID + ":" + strconv.Itoa(i.TryNumber))
		if sTask == nil {
			i.StartDate = time.Now()
			e.Redis.SaveDagTaskRedis(i)
			logrus.WithField("func", "EventLoop").Debug("Dag ID: ", i.DagID)
			logrus.WithField("func", "EventLoop").Debug("Dag Task ID: ", i.TaskID)
			logrus.WithField("func", "EventLoop").Debug("Dag Run ID: ", i.RunID)
			logrus.WithField("func", "EventLoop").Debug("Dag TryNumber: ", strconv.Itoa(i.TryNumber))
			logrus.WithField("func", "EventLoop").Debug("Dag StartDate: ", i.StartDate)
			logrus.WithField("func", "EventLoop").Debug("Dag CPUs: ", i.MesosExecutor.Cpus)
			logrus.WithField("func", "EventLoop").Debug("Dag MEM: ", i.MesosExecutor.MemLimit)
			logrus.WithField("func", "EventLoop").Debug("Dag InstanceType: ", i.MesosExecutor.InstanceType)
			logrus.WithField("func", "EventLoop").Debug("ASG: ", i.ASG)
			logrus.WithField("func", "EventLoop").Debug("---------------------------------------")
		}
	}
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

			var agent cfg.MesosAgentState
			err = json.NewDecoder(res.Body).Decode(&agent)
			if err != nil {
				logrus.WithField("func", "checkEC2Instances").Error("Could not encode json result: ", err.Error())
				// uncheck these instance
				instance.Check = false
				e.Redis.SaveEC2InstanceRedis(*instance)
				continue
			}

			// uncheck these instance
			instance.Check = false
			e.Redis.SaveEC2InstanceRedis(*instance)

			// check if there is still a airflow task running
			if !e.isFrameworkName(agent) {
				e.AWS.TerminateInstance(instance.EC2.Instances[0].InstanceId)
				e.Redis.DelRedisKey(e.Config.RedisPrefix + ":ec2:" + *instance.EC2.Instances[0].InstanceId)
			}
		}
	}
	if i > 0 {
		logrus.WithField("func", "mesos.checkEC2Instance").Infof("There are %d instances in DB. %d of them are in check mode.", i, c)
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
