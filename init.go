package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	util "github.com/AVENTER-UG/util"
	"github.com/sirupsen/logrus"

	cfg "github.com/AVENTER-UG/mesos-autoscale/types"
)

var config cfg.Config

func init() {

	config.AppName = "AWS Autoscale for Apache Airflow"
	config.AWSWait = util.Getenv("AWS_WAIT_TIME", "10m")
	config.AWSLaunchTemplateID = os.Getenv("AWS_LAUNCH_TEMPLATE_ID")
	config.AWSRegion = util.Getenv("AWS_REGION", "eu-central-1")
	config.AWSInstanceDefaultArchitecture = util.Getenv("AWS_INSTANCE_DEFAULT_ARCHITECTURE", "x86_64")
	config.AWSInstanceFallback = util.Getenv("AWS_INSTANCE_FALLBACK", "t3a.2xlarge")
	config.AirflowMesosScheduler = util.Getenv("AIRFLOW_MESOS_SCHEDULER", "127.0.0.1:11000")
	config.AirflowMesosName = util.Getenv("AIRFLOW_MESOS_NAME", "Airflow")
	config.APIUsername = util.Getenv("API_USERNAME", "user")
	config.APIPassword = util.Getenv("API_PASSWORD", "password")
	config.LogLevel = util.Getenv("LOGLEVEL", "debug")
	config.MesosAgentUsername = util.Getenv("MESOS_AGENT_USERNAME", "mesos")
	config.MesosAgentPassword = util.Getenv("MESOS_AGENT_PASSWORD", "")
	config.MesosAgentPort = util.Getenv("MESOS_AGENT_PORT", "5051")
	config.MesosMasterUsername = util.Getenv("MESOS_MASTER_USERNAME", "mesos")
	config.MesosMasterPassword = util.Getenv("MESOS_MASTER_PASSWORD", "")
	config.MesosMasterPort = util.Getenv("MESOS_MASTER_PORT", "5050")
	config.MesosMaster = util.Getenv("MESOS_MASTER", "leader.mesos")
	config.RedisServer = util.Getenv("REDIS_SERVER", "127.0.0.1:6480")
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	config.RedisDB, _ = strconv.Atoi(util.Getenv("REDIS_DB", "2"))
	config.RedisPrefix = util.Getenv("REDIS_PREFIX", "asg")

	// set pollinterval
	config.PollInterval, _ = time.ParseDuration(util.Getenv("POLL_INTERVALL", "2s"))
	config.PollTimeout, _ = time.ParseDuration(util.Getenv("POLL_TIMEOUT", "10s"))

	// set the time to wait that the dag is running until we scale up AWS
	config.WaitTimeout, _ = time.ParseDuration(util.Getenv("WAIT_TIME", "2m"))

	// set the time to wait until we overwrite the instance type from the executor.
	config.WaitTimeoutOverwrite, _ = time.ParseDuration(util.Getenv("WAIT_TIME_OVERWRITE_INSTANCE", "15m"))

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

	if strings.Compare(util.Getenv("AWS_INSTANCE_TERMINATE", "true"), "false") == 0 {
		config.AWSInstanceTerminate = false
	} else {
		config.AWSInstanceTerminate = true
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

	instanceTypes := strings.ReplaceAll(os.Getenv("AWS_INSTANCE_ALLOW"), "'", "\"")
	if instanceTypes != "" {
		err := json.Unmarshal([]byte(instanceTypes), &config.AWSInstanceAllow)

		if err != nil {
			logrus.WithField("func", "init").Fatal("The env variable AWS_INSTANCE_ALLOW have a syntax failure: ", err)
		}
	}
}
