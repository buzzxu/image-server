package local

import (
	"github.com/go-redis/redis/v7"
	"image-server/pkg/conf"
	"log"
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
	if _, err := cache.Ping().Result(); err != nil {
		log.Fatalf("Redis connect error.%s", err.Error())
	}
	log.Printf("Redis connect success")
}

func redisStats() {
	poolStats := cache.PoolStats()
	log.Printf("Redis Stats:[TotalConns:%d,IdleConns:%d,StaleConns:%d,Hits:%d,Misses:%d]",
		poolStats.TotalConns,
		poolStats.IdleConns,
		poolStats.StaleConns,
		poolStats.Hits,
		poolStats.Misses)
}
