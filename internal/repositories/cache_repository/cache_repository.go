package cache_repository

import (
	"context"
	"time"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value any, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}
