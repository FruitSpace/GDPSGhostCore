package api

import (
	"HalogenGhostCore/core"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var cacheClient *redis.Client

func InitCache(conf core.GlobalConfig, cacheDB int) {
	cacheClient = redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost + ":" + conf.RedisPort,
		Password: conf.RedisPassword,
		DB:       cacheDB,
	})
}

func cached(key string) (string, error) {
	ctx := context.Background()
	return cacheClient.Get(ctx, key).Result()
}

func withCache(key string, value string) string {
	ctx := context.Background()
	go cacheClient.Set(ctx, key, value, time.Second*2)
	return value
}
