package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(addr string) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisService{
		Client: client,
	}
}

func (rs *RedisService) Ping(ctx context.Context) error {
	return rs.Client.Ping(ctx).Err()
}

func (rs *RedisService) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return rs.Client.Set(ctx, key, value, exp).Err()
}

func (rs *RedisService) Get(ctx context.Context, key string) (string, error) {
	return rs.Client.Get(ctx, key).Result()
}

func (rs *RedisService) Del(ctx context.Context, key string) error {
	return rs.Client.Del(ctx, key).Err()
}
