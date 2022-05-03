package cache

import (
	"github.com/go-redis/redis"
	logger "github.com/sirupsen/logrus"
)

func InitRedis(host string, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
	})
	_, err := client.Ping().Result()
	if err != nil {
		logger.Fatal("Cann't connect to redis", err)
	}

	return client
}
