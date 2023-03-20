package task

import (
	"chat/config"
	"chat/tools"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func (task *Task) InitQueueRedisClient() (err error) {
	redisOpt := tools.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	RedisClient = tools.GetRedisInstance(redisOpt)

	if pong, err := RedisClient.Ping().Result(); err != nil {
		logrus.Infof("RedisClient Ping Result pong: %s,  err: %s", pong, err)
	}

	//开启协程，循环等待api层的push
	go func() {
		for {
			var result []string
			//读取数据，取数据时如果没有数据会等待指定的时间
			result, err = RedisClient.BRPop(time.Second*10, config.QueueName).Result()
			if err != nil {
				logrus.Infof("task queue block timeout,no msg err:%s", err.Error())
			}
			//BRPOP读到的数据会是长度为2的数组，其中result[0]为键，result[1]为值
			if len(result) >= 2 {
				task.Push(result[1])
			}
		}
	}()
	return
}
