package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"

	"github.com/mohit838/olario-platform-backend/internal/config"
)

// Open creates a Redis client and verifies the connection.
// Redis is used for fast counters/cache-style data; durable business state
// stays in Postgres.
func Open(ctx context.Context, cfg config.RedisConfig) (*goredis.Client, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("redis addr is required")
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
