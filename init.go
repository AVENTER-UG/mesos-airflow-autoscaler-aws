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
	config.AirflowMesosName = util.Getenv("AIRFLOW_MESOS_NAME", "Airflow")
	config.LogLevel = util.Getenv("LOGLEVEL", "debug")
	config.AWSWait = util.Getenv("AWS_WAIT_TIME", "10m")
	config.AppName = "AWS Autoscale for Apache Airflow"
	config.RedisServer = util.Getenv("REDIS_SERVER", "127.0.0.1:6480")
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	config.RedisDB, _ = strconv.Atoi(util.Getenv("REDIS_DB", "2"))
	config.RedisPrefix = util.Getenv("REDIS_PREFIX", "asg")
	config.APIUsername = util.Getenv("API_USERNAME", "user")
	config.APIPassword = util.Getenv("API_PASSWORD", "password")
	config.AWSLaunchTemplateID = os.Getenv("AWS_LAUNCH_TEMPLATE_ID")
	config.AWSRegion = util.Getenv("AWS_REGION", "eu-central-1")
	config.AWSInstance16 = util.Getenv("AWS_INSTANCE_16GB", "t2.xlarge")
	config.AWSInstance32 = util.Getenv("AWS_INSTANCE_32GB", "t3.2xlarge")
	config.AWSInstance64 = util.Getenv("AWS_INSTANCE_64GB", "r5.2xlarge")
	config.MesosAgentUsername = os.Getenv("MESOS_AGENT_USERNAME")
	config.MesosAgentPassword = os.Getenv("MESOS_AGENT_PASSWORD")
	config.MesosAgentPort = util.Getenv("MESOS_AGENT_PORT", "5051")

	// set pollinterval
	config.PollInterval, _ = time.ParseDuration(util.Getenv("POLL_INTERVALL", "2s"))
	config.PollTimeout, _ = time.ParseDuration(util.Getenv("POLL_TIMEOUT", "10s"))

	// set the time to wait that the dag is running until we scale up AWS
	config.WaitTimeout, _ = time.ParseDuration(util.Getenv("WAIT_TIME", "2m"))

	// set the time to wait until we will check if we can terminte the ec2 instance
	config.AWSTerminateWait, _ = time.ParseDuration(config.AWSWait)

	// set mesos agent timeout
	config.MesosAgentTimeout, _ = time.ParseDuration(util.Getenv("MESOS_AGENT_TIMEOUT", "10m"))

	// set TTL for dags in redis
	config.RedisTTL, _ = time.ParseDuration(util.Getenv("DAG_TTL", "6h"))

	if strings.Compare(util.Getenv("SSL", "false"), "true") == 0 {
		config.SSL = true
	} else {
		config.SSL = false
	}

	// The comunication to the mesos server should be via ssl or not
	if strings.Compare(util.Getenv("MESOS_AGENT_SSL", "false"), "true") == 0 {
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
	if strings.Compare(util.Getenv("SKIP_SSL", "true"), "true") == 0 {
		config.SkipSSL = true
	} else {
		config.SkipSSL = false
	}
}
