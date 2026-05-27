package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// DevCache is the Redis adapter used by the local full-circle demo API.
type DevCache struct {
	client *goredis.Client
}

func NewDevCache(client *goredis.Client) *DevCache {
	return &DevCache{client: client}
}

func (c *DevCache) IncrementFullCircleRuns(ctx context.Context) (int64, error) {
	return c.client.Incr(ctx, "dev:full_circle:runs").Result()
}

func (c *DevCache) RememberLastOrder(ctx context.Context, orderNumber string) error {
	return c.client.Set(ctx, "dev:full_circle:last_order", orderNumber, 24*time.Hour).Err()
}
