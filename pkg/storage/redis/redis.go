package redis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"image-server/pkg/conf"
)

// Redis 客户端
var Redis *redis.Client

// RedisConnect Redis连接
func RedisConnect() {
	var password = ""
	var poolSize = 100
	if len(conf.Config.Redis.Password) > 0 {
		password = conf.Config.Redis.Password
	}
	if conf.Config.Redis.PoolSize > 0 {
		poolSize = conf.Config.Redis.PoolSize
	}
	Redis = redis.NewClient(&redis.Options{
		Addr:     conf.Config.Redis.Addr,
		Password: password,             // no password set
		DB:       conf.Config.Redis.DB, // use default DB
		PoolSize: poolSize,
	})

	pong, err := Redis.Ping().Result()
	fmt.Println(pong, err)
}
