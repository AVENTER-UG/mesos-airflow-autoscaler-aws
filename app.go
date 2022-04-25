package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	goredis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// BuildVersion of m3s
var BuildVersion string

// GitVersion is the revision and commit number
var GitVersion string

// init the redis cache
func initCache() {
	var redisOptions goredis.Options
	redisOptions.Addr = config.RedisServer
	redisOptions.DB = config.RedisDB
	if config.RedisPassword != "" {
		redisOptions.Password = config.RedisPassword
	}
	client := goredis.NewClient(&redisOptions)

	config.RedisCTX = context.Background()
	pong, err := client.Ping(config.RedisCTX).Result()
	logrus.Debug("Redis Health: ", pong, err)
	config.RedisClient = client
}

// convert Base64 Encodes PEM Certificate to tls object
func decodeBase64Cert(pemCert string) []byte {
	sslPem, err := base64.URLEncoding.DecodeString(pemCert)
	if err != nil {
		logrus.Fatal("Error decoding SSL PEM from Base64: ", err.Error())
	}
	return sslPem
}

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

	initCache()

	logrus.Fatal(subscribe())
}
