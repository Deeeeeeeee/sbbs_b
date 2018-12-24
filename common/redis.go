package common

import (
	"sbbs_b/config"
	"sync"

	"github.com/go-redis/redis"
)

var (
	client    *redis.Client
	redisOnce sync.Once
)

// Redis 返回 redis client
func Redis() *redis.Client {
	redisLazyinit()
	return client
}

func redisLazyinit() {
	redisOnce.Do(func() {
		initClient()
	})
}

func initClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
}
