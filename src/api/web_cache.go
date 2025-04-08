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
	return withCacheDuration(key, value, time.Second*2)
}

func withCacheDuration(key string, value string, duration time.Duration) string {
	ctx := context.Background()
	cacheClient.Set(ctx, key, value, duration)
	return value
}
