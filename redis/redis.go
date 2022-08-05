package redis

import (
	"context"
	"encoding/json"
	"strconv"

	cfg "github.com/AVENTER-UG/mesos-autoscale/types"
	"github.com/aws/aws-sdk-go/service/ec2"
	goredis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// Redis struct about the redis connection
type Redis struct {
	RedisClient *goredis.Client
	RedisCTX    context.Context
	Config      *cfg.Config
}

// New will create a new API object
func New(cfg *cfg.Config) *Redis {
	e := &Redis{
		Config: cfg,
	}

	return e
}

// GetAllRedisKeys get out all keys in redis depends to the pattern
func (e *Redis) GetAllRedisKeys(pattern string) *goredis.ScanIterator {
	val := e.RedisClient.Scan(e.RedisCTX, 0, pattern, 0).Iterator()
	if err := val.Err(); err != nil {
		logrus.Warn("getAllRedisKeys: ", err.Error())
	}
	return val
}

// GetRedisKey get out all values to a key
func (e *Redis) GetRedisKey(key string) string {
	val, err := e.RedisClient.Get(e.RedisCTX, key).Result()
	if err != nil {
		logrus.Warn("getRedisKey: ", err.Error())
	}
	return val
}

// DelRedisKey will delete a redis key
func (e *Redis) DelRedisKey(key string) int64 {
	val, err := e.RedisClient.Del(e.RedisCTX, key).Result()
	if err != nil {
		logrus.Warn("delRedisKey: ", err.Error())
	}

	return val
}

// GetTaskFromRunID get out the task to an runID
func (e *Redis) GetTaskFromRunID(key string) *cfg.DagTask {
	// search matched taskid in redis and update the status
	keys := e.GetRedisKey(key)
	if len(keys) > 0 {
		var task cfg.DagTask
		json.Unmarshal([]byte(keys), &task)
		return &task
	}

	return nil
}

// GetEC2InstanceFromID get out the task to an runID
func (e *Redis) GetEC2InstanceFromID(key string) *ec2.Reservation {
	// search matched taskid in redis and update the status
	keys := e.GetRedisKey(key)
	if len(keys) > 0 {
		var instance *ec2.Reservation
		json.Unmarshal([]byte(keys), &instance)
		return instance
	}

	return nil
}

// CountRedisKey will get back the count of the redis key
func (e *Redis) CountRedisKey(pattern string) int {
	keys := e.GetAllRedisKeys(pattern)
	count := 0
	for keys.Next(e.RedisCTX) {
		count++
	}
	logrus.Debug("CountRedisKey: ", pattern, count)
	return count
}

// SaveTaskRedis store mesos task in DB
func (e *Redis) SaveDagTaskRedis(task cfg.DagTask) {
	d, _ := json.Marshal(&task)
	e.SetRedisKey(d, task.DagID+":"+task.TaskID+":"+task.RunID+":"+strconv.Itoa(task.TryNumber))
}

// SaveEC2InstanceRedis store mesos task in DB
func (e *Redis) SaveEC2InstanceRedis(instance *ec2.Reservation) {
	d, _ := json.Marshal(&instance)
	err := e.RedisClient.Set(e.RedisCTX, e.Config.RedisPrefix+":ec2:"+*instance.Instances[0].InstanceId, d, 0).Err()
	if err != nil {
		logrus.WithField("func", "SaveData").Error("Could not save data in Redis: ", err.Error())
	}
}

// SetRedisKey store data in redis
func (e *Redis) SetRedisKey(data []byte, key string) {
	err := e.RedisClient.Set(e.RedisCTX, e.Config.RedisPrefix+":dags:"+key, data, 0).Err()
	if err != nil {
		logrus.WithField("func", "SaveData").Error("Could not save data in Redis: ", err.Error())
	}
}

// PingRedis to check the health of redis
func (e *Redis) PingRedis() error {
	pong, err := e.RedisClient.Ping(e.RedisCTX).Result()
	logrus.Debug("Redis Health: ", pong, err)
	if err != nil {
		return err
	}
	return nil
}

// ConnectRedis will connect the redis DB and save the client pointer
func (e *Redis) ConnectRedis() {
	var redisOptions goredis.Options
	redisOptions.Addr = e.Config.RedisServer
	redisOptions.DB = e.Config.RedisDB
	if e.Config.RedisPassword != "" {
		redisOptions.Password = e.Config.RedisPassword
	}

	e.RedisClient = goredis.NewClient(&redisOptions)
	e.RedisCTX = context.Background()

	err := e.PingRedis()
	if err != nil {
		e.ConnectRedis()
	}
}
