package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"lls_api/pkg/config"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

var r *RedisClient

func InitRedis() {
	r = &RedisClient{client: redis.NewClient(&redis.Options{
		Addr:     config.C.Redis.Host,
		Password: config.C.Redis.Pwd,
		DB:       config.C.Redis.DB,
	})}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}
}

func FullKey(key string) string {
	return config.C.Redis.Prefix + ":" + key
}

func (c *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	key = FullKey(key)
	return r.client.Get(ctx, key)
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	key = FullKey(key)
	return r.client.Set(ctx, key, value, expiration)
}

func (c *RedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	var fullKeys []string
	for _, e := range keys {
		fullKeys = append(fullKeys, FullKey(e))
	}
	return r.client.Del(ctx, fullKeys...)
}
