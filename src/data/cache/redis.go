package cache

import (
	"fmt"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/go-redis/redis/v8"
	"time"
)

var redisClient *redis.Client

func InitRedis(cfg *config.Config) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:           cfg.Redis.Password,
		DB:                 0,
		DialTimeout:        cfg.Redis.DialTimeout * time.Second,
		ReadTimeout:        cfg.Redis.ReadTimeout * time.Second,
		WriteTimeout:       cfg.Redis.WriteTimeout * time.Second,
		PoolSize:           cfg.Redis.PoolSize,
		PoolTimeout:        cfg.Redis.PoolTimeout * time.Second,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: cfg.Redis.IdleCheckFrequency * time.Second,
	})
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	redisClient.Close()
}