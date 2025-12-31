package cache_repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheRepository struct {
	Client *redis.Client
}

func NewRedisCacheRepository(addr string) *RedisCacheRepository {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisCacheRepository{
		Client: client,
	}
}

func (rs *RedisCacheRepository) Ping(ctx context.Context) error {
	return rs.Client.Ping(ctx).Err()
}

func (rs *RedisCacheRepository) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return rs.Client.Set(ctx, key, value, exp).Err()
}

func (rs *RedisCacheRepository) Get(ctx context.Context, key string) (string, error) {
	return rs.Client.Get(ctx, key).Result()
}

func (rs *RedisCacheRepository) Del(ctx context.Context, key string) error {
	return rs.Client.Del(ctx, key).Err()
}
