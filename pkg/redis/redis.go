package redis

import (
	"github.com/go-redis/redis/v7"
	"image-server/pkg/conf"
	"log"
	"time"
)

// Redis 客户端
var Client *redis.Client

var ticker *time.Ticker

// RedisConnect Redis连接
func RedisConnect() {
	Client = redis.NewClient(&redis.Options{
		Addr:         conf.Config.Redis.Addr,
		Password:     conf.Config.Redis.Password, // no password set
		DB:           conf.Config.Redis.DB,       // use default DB
		PoolSize:     conf.Config.Redis.PoolSize,
		MaxRetries:   3,
		MinIdleConns: conf.Config.Redis.MinIdleConns,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		IdleTimeout:  60 * time.Second,
	})
	if _, err := Client.Ping().Result(); err != nil {
		log.Fatalf("Redis connect error.%s", err.Error())
	}
	log.Printf("Redis connect success")

	ticker = time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				RedisStats()
			}
		}
	}()
}

func RedisStats() {
	poolStats := Client.PoolStats()
	log.Printf("Redis Stats:[TotalConns:%d,IdleConns:%d,StaleConns:%d,Hits:%d,Misses:%d]",
		poolStats.TotalConns,
		poolStats.IdleConns,
		poolStats.StaleConns,
		poolStats.Hits,
		poolStats.Misses)
}

func Close() {
	ticker.Stop()
	Client.Close()
}
