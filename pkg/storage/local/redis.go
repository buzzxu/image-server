package local

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"image-server/pkg/conf"
	"time"
)

// Redis 客户端
var cache *redis.Client

// RedisConnect Redis连接
func redisConnect() {
	cache = redis.NewClient(&redis.Options{
		Addr:         conf.Config.Redis.Addr,
		Password:     conf.Config.Redis.Password, // no password set
		DB:           conf.Config.Redis.DB,       // use default DB
		PoolSize:     conf.Config.Redis.PoolSize,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		PoolTimeout:  20 * time.Second,
	})
	pong, err := cache.Ping().Result()
	fmt.Println(pong, err)
}
