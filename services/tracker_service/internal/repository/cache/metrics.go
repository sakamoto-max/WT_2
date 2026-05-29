package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type metricsCache struct {
	client *redis.Client
}

func (c *metricsCache) GetRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := c.client.Ping(ctx).Err()
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
