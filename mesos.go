package main

import (
	"context"
	cTls "crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)


func subscribe() {
		// Add protocoll to the endpoint depends if SSL is enabled
		logrus.Info("Connect Mesos Provider to: ", config.MesosMasterServer)

		operation := func() error {
			ticker := time.NewTicker(time.Duration(PollInterval))
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					data, _ := getTasks()

					var tasks MesosTasks
					if err := json.Unmarshal(data, &tasks); err != nil {
						logger.Error("Error in Data from Mesos: " + err.Error())
						continue
					}

					// collect all mesos tasks and combine the belong one.
					for _, task := range tasks.Tasks {
						switch task.State {
						case "TASK_RUNNING":
						}
					}

					// cleanup old data
					mesosConfig = make(map[string]*MesosTasks)
				case <-routineCtx.Done():
					return nil
				}
			}
			return nil
		}
}

func getTasks() ([]byte, error) {
	client := &httClient{}
	client.Transport = &httTransport{
		TLSClientConfig: &cTls.Config{InsecureSkipVerify: true},
	}
	req, _ := httNewRequest("GET", config.MesosMasterServer+"/tasks?order=asc&limit=-1", nil)
	req.Close = true
	req.SetBasicAuth(config.Username, config.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != httStatusOK {
		return nil, fmt.Errorf("received non-ok response code: %d", res.StatusCode)
	}

	logger.Info("Get Data from Mesos")
	return io.ReadAll(res.Body)
}

func checkContainer(task MesosTask) bool {
	agentHostname, agentPort, err := getAgent(task.SlaveID)

	if err != nil {
		logger.Error("CheckContainer: Error in get AgendData from Mesos: " + err.Error())
		return false
	}

	logger.Debug("CheckContainer: " + task.Name + " on agent (" + task.SlaveID + ")" + agentHostname + " with task.ID " + task.ID)

	if agentHostname != "" {
		containers, _ := getContainersOfAgent(agentHostname, agentPort)

		for _, a := range containers {
			logger.Debug(task.ID + " --CONTAINER--  " + a.ExecutorID)
			if a.ExecutorID == task.ID {
				return true
			}
		}
	}

	return false
}

func getAgent(slaveID string) (string, int, error) {
	client := &httClient{}
	client.Transport = &httTransport{
		TLSClientConfig: &cTls.Config{InsecureSkipVerify: true},
	}
	req, _ := httNewRequest("GET", config.MesosMasterServer+"/slaves/", nil)
	req.Close = true
	req.SetBasicAuth(config.Username, config.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != httStatusOK {
		return "", 0, fmt.Errorf("received non-ok response code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	var agents MesosAgent
	if err := json.Unmarshal(data, &agents); err != nil {
		logger.Error("getAgent: Error in AgentData from Mesos: " + err.Error())
		return "", 0, err
	}

	for _, a := range agents.Slaves {
		if a.ID == slaveID {
			return a.Hostname, a.Port, nil
		}
	}

	return "", 0, nil
}

func getContainersOfAgent(agentHostname string, agentPort int) (MesosAgentContainers, error) {
	// Add protocoll to the endpoint depends if SSL is enabled
	protocol := "http://"
	if SSL {
		protocol = "https://"
	}

	client := &httClient{}
	client.Transport = &httTransport{
		TLSClientConfig: &cTls.Config{InsecureSkipVerify: true},
	}
	req, _ := httNewRequest("GET", protocol+agentHostname+":"+strconv.Itoa(agentPort)+"/containers/", nil)
	req.Close = true
	req.SetBasicAuth(Principal, Secret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != httStatusOK {
		return MesosAgentContainers{}, fmt.Errorf("received non-ok response code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	var containers MesosAgentContainers
	if err := json.Unmarshal(data, &containers); err != nil {
		logger.Error("getContainersOfAgent: Error in ContainerAgentData from " + agentHostname + "  " + err.Error())
		return MesosAgentContainers{}, err
	}

	return containers, nil
}
