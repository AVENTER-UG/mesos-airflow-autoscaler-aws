package main

import (
	"os"
	"strconv"
	"strings"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	util "github.com/AVENTER-UG/util"
	"github.com/Showmax/go-fqdn"

	cfg "github.com/AVENTER-UG/mesos-autoscale/types"
)

var config cfg.Config

func init() {
	config.Username = os.Getenv("MESOS_USERNAME")
	config.Password = os.Getenv("MESOS_PASSWORD")
	config.MesosMasterServer = os.Getenv("MESOS_MASTER")
	config.LogLevel = util.Getenv("LOGLEVEL", "info")
	config.AppName = "Mesos Autoscale"
	config.RedisServer = util.Getenv("REDIS_SERVER", "127.0.0.1:6379")
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	config.RedisDB, _ = strconv.Atoi(util.Getenv("REDIS_DB", "1"))

	// The comunication to the mesos server should be via ssl or not
	if strings.Compare(os.Getenv("MESOS_SSL"), "true") == 0 {
		config.MesosSSL = true
	} else {
		config.MesosSSL = false
	}

  protocol := "http://" + config.MesosMasterServer
  if config.MesosSSL  {
  	protocol = "https://" + config.MesosMasterServer
  }
	config.MesosMasterServer = protocol

	// Skip SSL Verification
	if strings.Compare(os.Getenv("SKIP_SSL"), "true") == 0 {
		config.SkipSSL = true
	} else {
		config.SkipSSL = false
	}

	config.PollInterval = ptypes.Duration(10 * time.Second)
	config.PollTimeout = ptypes.Duration(10 * time.Second)
}
