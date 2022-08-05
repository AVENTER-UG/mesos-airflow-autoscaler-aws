package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	util "github.com/AVENTER-UG/util"

	cfg "github.com/AVENTER-UG/mesos-autoscale/types"
)

var config cfg.Config

func init() {

	config.AirflowMesosScheduler = util.Getenv("AIRFLOW_MESOS_SCHEDULER", "127.0.0.1:11000")
	config.LogLevel = util.Getenv("LOGLEVEL", "debug")
	config.Wait = util.Getenv("WAIT_TIME", "2m")
	config.AWSWait = util.Getenv("AWS_WAIT_TIME", "5m")
	config.AppName = "AWS Autoscale for Apache Airflow"
	config.RedisServer = util.Getenv("REDIS_SERVER", "127.0.0.1:6480")
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	config.RedisDB, _ = strconv.Atoi(util.Getenv("REDIS_DB", "2"))
	config.RedisPrefix = util.Getenv("REDIS_PREFIX", "asg")
	config.APIUsername = util.Getenv("API_USERNAME", "user")
	config.APIPassword = util.Getenv("API_PASSWORD", "password")
	config.AWSLaunchTemplateID = os.Getenv("AWS_LAUNCH_TEMPLATE_ID")
	config.AWSRegion = util.Getenv("AWS_REGION", "eu-central-1")
	config.MesosAgentUsername = os.Getenv("MESOS_AGENT_USERNAME")
	config.MesosAgentPassword = os.Getenv("MESOS_AGENT_PASSWORD")
	config.MesosAgentPort = util.Getenv("MESOS_AGENT_PORT", "5051")
	config.PollInterval = 5 * time.Second
	config.PollTimeout = 10 * time.Second

	// set the time to wait that the dag is running until we scale up AWS
	config.WaitTimeout, _ = time.ParseDuration(config.Wait)

	// set the time to wait until we will check if we can terminte the ec2 instance
	config.AWSTerminateWait, _ = time.ParseDuration(config.AWSWait)

	if strings.Compare(os.Getenv("SSL"), "true") == 0 {
		config.SSL = true
	} else {
		config.SSL = false
	}

	// The comunication to the mesos server should be via ssl or not
	if strings.Compare(os.Getenv("MESOS_AGENT_SSL"), "true") == 0 {
		config.MesosAgentSSL = true
	} else {
		config.MesosAgentSSL = false
	}

	protocol := "http://" + config.AirflowMesosScheduler
	if config.SSL {
		protocol = "https://" + config.AirflowMesosScheduler
	}

	config.AirflowMesosScheduler = protocol

	// Skip SSL Verification
	if strings.Compare(os.Getenv("SKIP_SSL"), "true") == 0 {
		config.SkipSSL = true
	} else {
		config.SkipSSL = false
	}
}
