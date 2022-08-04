package types

import (
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// Config hold the configuration of these software
type Config struct {
	MinVersion            string
	EnableSyslog          bool
	AirflowMesosScheduler string
	LogLevel              string
	AppName               string
	Wait                  string
	SSL                   bool
	SkipSSL               bool
	PollInterval          time.Duration
	PollTimeout           time.Duration
	WaitTimeout           time.Duration
	RedisServer           string
	RedisPassword         string
	RedisDB               int
	RedisPrefix           string
	APIUsername           string
	APIPassword           string
	AWSRegion             string
	AWSLaunchTemplateID   string
}

type DagTask struct {
	DagID         string `json:"dag_id"`
	TaskID        string `json:"task_id"`
	RunID         string `json:"run_id"`
	TryNumber     int    `json:"try_number"`
	ASG           bool   `json:"ASG" default:"false"`
	EC2           *ec2.Reservation
	MesosExecutor struct {
		Cpus     float64 `json:"cpus"`
		MemLimit int     `json:"mem_limit"`
	} `json:"MesosExecutor"`
}
