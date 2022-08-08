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
}

// New will create the Scheduler object
func New(cfg *cfg.Config) *Scheduler {
	e := &Scheduler{
		Config: cfg,
	}
	// Add protocoll to the endpoint depends if SSL is enabled
	logrus.Info("Connect Provider Apache Mesos Airflow Provider: ", e.Config.AirflowMesosScheduler)

	return e
}

// EventLoop is the main loop of all events
func (e *Scheduler) EventLoop() {
	ticker := time.NewTicker(e.Config.PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			data := e.getDags()

			for _, i := range data {
				// get execution time of dag task

				timeStr := strings.Split(i.RunID, "__")
				timeTask, err := time.Parse(time.RFC3339, timeStr[1])

				// get current time
				timeNow := time.Now()
				timeDiff := timeNow.Sub(timeTask).Seconds()

				if err != nil {
					logrus.WithField("func", "EventLoop").Error("Cannot parse TimeStamp: ", err.Error())
					continue
				}
				logrus.WithField("func", "EventLoop").Debug("Dag ID: ", i.DagID)
				logrus.WithField("func", "EventLoop").Debug("Dag Task ID: ", i.TaskID)
				logrus.WithField("func", "EventLoop").Debug("Dag Run ID: ", i.RunID)
				logrus.WithField("func", "EventLoop").Debug("Dag TryNumber: ", strconv.Itoa(i.TryNumber))
				logrus.WithField("func", "EventLoop").Debug("Dag Task Age: ", timeDiff)
				logrus.WithField("func", "EventLoop").Debug("Dag CPUs: ", i.MesosExecutor.Cpus)
				logrus.WithField("func", "EventLoop").Debug("Dag MEM: ", i.MesosExecutor.MemLimit)
				logrus.WithField("func", "EventLoop").Debug("ASG: ", i.ASG)
				logrus.WithField("func", "EventLoop").Debug("---------------------------------------")

				// check if the runID already exist
				sTask := e.Redis.GetTaskFromRunID(e.Config.RedisPrefix + ":dags:" + i.DagID + ":" + i.TaskID + ":" + i.RunID + ":" + strconv.Itoa(i.TryNumber))
				if sTask != nil {
					i = *sTask
					logrus.WithField("func", "EventLoop").Debug("=== ASG Already ")
				} else {
					if timeDiff >= e.Config.WaitTimeout.Seconds() {
						logrus.WithField("func", "EventLoop").Debug(">>> ASG ScaleUp ")
						i.ASG = true

						if i.MesosExecutor.MemLimit >= 32768 {
							go e.Redis.SaveEC2InstanceRedis(e.AWS.CreateInstance(e.Config.AWSInstance64))
						} else if i.MesosExecutor.MemLimit >= 16384 {
							go e.Redis.SaveEC2InstanceRedis(e.AWS.CreateInstance(e.Config.AWSInstance32))
						} else {
							go e.Redis.SaveEC2InstanceRedis(e.AWS.CreateInstance(e.Config.AWSInstance16))
						}
						e.Redis.SaveDagTaskRedis(i)
					}
				}

			}
			// Check if we can terminate the ec2 instances
			go e.checkEC2Instance()
		}
	}
}

// get all running dags from airflow mesos scheduler
func (e *Scheduler) getDags() []cfg.DagTask {
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
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logrus.WithField("func", "getDags").Error("Status Code not 200: ", res.StatusCode)
		return nil
	}

	logrus.Info("Get Data from Mesos")
	var dags []cfg.DagTask
	err = json.NewDecoder(res.Body).Decode(&dags)
	if err != nil {
		logrus.WithField("func", "getDags").Error("Cannot decode json: ", err.Error())
		return nil
	}
	return dags
}

// check if the ec2 instance still running mesos tasks
func (e *Scheduler) checkEC2Instance() {
	logrus.WithField("func", "checkEC2Instance").Info("Check EC2 Instances")
	keys := e.Redis.GetAllRedisKeys(e.Config.RedisPrefix + ":ec2:*")
	for keys.Next(e.Redis.RedisCTX) {
		logrus.WithField("func", "checkEC2Instance").Debug("Instances: ", keys.Val())
		instance := e.Redis.GetEC2InstanceFromID(keys.Val())

		// if the launch time is to short, do not try to connect reach the agent
		// get current time
		timeNow := time.Now()
		timeDiff := timeNow.Sub(*instance.Instances[0].LaunchTime).Minutes()

		if timeDiff <= e.Config.AWSTerminateWait.Minutes() {
			continue
		}

		client := &http.Client{
			Timeout: 2 * time.Second,
		}
		client.Transport = &http.Transport{
			// #nosec G402
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		protocol := "https"
		if !e.Config.MesosAgentSSL {
			protocol = "http"
		}

		hostIP := instance.Instances[0].NetworkInterfaces[0].PrivateIpAddress

		req, _ := http.NewRequest("POST", protocol+"://"+*hostIP+":"+e.Config.MesosAgentPort+"/state", nil)
		req.Close = true
		req.SetBasicAuth(e.Config.MesosAgentUsername, e.Config.MesosAgentPassword)
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)

		if err != nil {
			logrus.WithField("func", "checkEC2Instances").Error("Could not connect to agent: ", err.Error())
			e.AWS.TerminateInstance(instance.Instances[0].InstanceId)
			e.Redis.DelRedisKey(e.Config.RedisPrefix + ":ec2:" + *instance.Instances[0].InstanceId)
			continue
		}

		defer res.Body.Close()

		var agent cfg.MesosAgentState
		err = json.NewDecoder(res.Body).Decode(&agent)
		if err != nil {
			logrus.WithField("func", "checkEC2Instances").Error("Could not encode json result: ", err.Error())
			continue
		}

		if len(agent.Frameworks) <= 0 {
			e.AWS.TerminateInstance(instance.Instances[0].InstanceId)
			e.Redis.DelRedisKey(e.Config.RedisPrefix + ":ec2:" + *instance.Instances[0].InstanceId)
		}
	}
}
