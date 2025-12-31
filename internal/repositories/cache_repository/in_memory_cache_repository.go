package cache_repository

import (
	"context"
	"time"
)

type InMemoryCacheRepository struct {
	data map[string]string
}

func NewInMemoryCacheRepository() *InMemoryCacheRepository {
	return &InMemoryCacheRepository{
		data: make(map[string]string),
	}
}

func (r *InMemoryCacheRepository) Set(ctx context.Context, key string, value any, exp time.Duration) error {
	r.data[key] = value.(string)
	return nil
}

func (r *InMemoryCacheRepository) Get(ctx context.Context, key string) (string, error) {
	value, exists := r.data[key]
	if !exists {
		return "", nil
	}
	return value, nil
}

func (r *InMemoryCacheRepository) Del(ctx context.Context, key string) error {
	delete(r.data, key)
	return nil
}
