package types

import (
	"time"
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
	AWSSecret             string
	AWSRegion             string
	AWSImageID            string
}

type DagTask struct {
	DagID         string `json:"dag_id"`
	TaskID        string `json:"task_id"`
	RunID         string `json:"run_id"`
	TryNumber     int    `json:"try_number"`
	ASG           bool   `json:"ASG" default:"false"`
	MesosExecutor struct {
		Cpus     float64 `json:"cpus"`
		MemLimit int     `json:"mem_limit"`
	} `json:"MesosExecutor"`
}
