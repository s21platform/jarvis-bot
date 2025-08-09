package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/s21platform/jarvis-bot/internal/config"
)

type Repository struct {
	client *redis.Client
}

func New(cfg *config.Config) *Repository {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     "",
		DB:           0,
		MinIdleConns: 2,
	})
	return &Repository{client: client}
}

func (r *Repository) Close() error {
	return r.client.Close()
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("cannot get value by key: %s, err: %v", key, err)
	}
	return res, nil
}

func (r *Repository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}
