package mesos

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

	"github.com/cenkalti/backoff/v4"
	ptypes "github.com/traefik/paerser/types"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/job"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/provider"
	"github.com/traefik/traefik/v2/pkg/safe"

	// Register mesos zoo the detector

	_ "github.com/mesos/mesos-go/api/v0/detector/zoo"
)


func connectMesos() {
		// Add protocoll to the endpoint depends if SSL is enabled
		protocol := "http://" + p.Endpoint
		if p.SSL {
			protocol = "https://" + p.Endpoint
		}
		p.Endpoint = protocol

		p.logger.Info("Connect Mesos Provider to: ", p.Endpoint)

		operation := func() error {
			ticker := time.NewTicker(time.Duration(p.PollInterval))
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					data, _ := getTasks()

					var tasks MesosTasks
					if err := json.Unmarshal(data, &tasks); err != nil {
						p.logger.Error("Error in Data from Mesos: " + err.Error())
						continue
					}

					// collect all mesos tasks and combine the belong one.
					for _, task := range tasks.Tasks {
						switch task.State {
						case "TASK_RUNNING":
						}
					}

					// cleanup old data
					p.mesosConfig = make(map[string]*MesosTasks)
				case <-routineCtx.Done():
					return nil
				}
			}
			return nil
		}
}

func getTasks() ([]byte, error) {
	client := &http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &cTls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("GET", p.Endpoint+"/tasks?order=asc&limit=-1", nil)
	req.Close = true
	req.SetBasicAuth(p.Principal, p.Secret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		p.logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-ok response code: %d", res.StatusCode)
	}

	p.logger.Info("Get Data from Mesos")
	return io.ReadAll(res.Body)
}

func checkContainer(task MesosTask) bool {
	agentHostname, agentPort, err := p.getAgent(task.SlaveID)

	if err != nil {
		p.logger.Error("CheckContainer: Error in get AgendData from Mesos: " + err.Error())
		return false
	}

	p.logger.Debug("CheckContainer: " + task.Name + " on agent (" + task.SlaveID + ")" + agentHostname + " with task.ID " + task.ID)

	if agentHostname != "" {
		containers, _ := p.getContainersOfAgent(agentHostname, agentPort)

		for _, a := range containers {
			p.logger.Debug(task.ID + " --CONTAINER--  " + a.ExecutorID)
			if a.ExecutorID == task.ID {
				return true
			}
		}
	}

	return false
}

func getAgent(slaveID string) (string, int, error) {
	client := &http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &cTls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("GET", p.Endpoint+"/slaves/", nil)
	req.Close = true
	req.SetBasicAuth(p.Principal, p.Secret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		p.logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("received non-ok response code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	var agents MesosAgent
	if err := json.Unmarshal(data, &agents); err != nil {
		p.logger.Error("getAgent: Error in AgentData from Mesos: " + err.Error())
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
	if p.SSL {
		protocol = "https://"
	}

	client := &http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &cTls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("GET", protocol+agentHostname+":"+strconv.Itoa(agentPort)+"/containers/", nil)
	req.Close = true
	req.SetBasicAuth(p.Principal, p.Secret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		p.logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return MesosAgentContainers{}, fmt.Errorf("received non-ok response code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	var containers MesosAgentContainers
	if err := json.Unmarshal(data, &containers); err != nil {
		p.logger.Error("getContainersOfAgent: Error in ContainerAgentData from " + agentHostname + "  " + err.Error())
		return MesosAgentContainers{}, err
	}

	return containers, nil
}
