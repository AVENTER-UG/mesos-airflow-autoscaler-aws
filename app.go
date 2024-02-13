package main

import (
	"flag"
	"fmt"

	mesosaws "github.com/AVENTER-UG/mesos-autoscale/aws"
	"github.com/AVENTER-UG/mesos-autoscale/mesos"
	"github.com/AVENTER-UG/mesos-autoscale/redis"
	util "github.com/AVENTER-UG/util/util"
	"github.com/sirupsen/logrus"
)

// BuildVersion of m3s
var BuildVersion string

// GitVersion is the revision and commit number
var GitVersion string

func main() {
	// Prints out current version
	var version bool
	flag.BoolVar(&version, "v", false, "Prints current version")
	flag.Parse()
	if version {
		fmt.Print(GitVersion)
		return
	}

	util.SetLogging(config.LogLevel, config.EnableSyslog, config.AppName)
	logrus.Println(config.AppName + " build " + BuildVersion + " git " + GitVersion)

	// connect to redis db
	r := redis.New(&config)
	r.ConnectRedis()

	// connect to aws
	a := mesosaws.New(&config)

	e := mesos.New(&config)
	e.Redis = r
	e.AWS = a
	go e.HealthCheck()
	e.EventLoop()
}
