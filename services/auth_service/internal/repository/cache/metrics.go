package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type metrics struct {
	client *redis.Client
}



func (c *metrics) GetRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()

	err := c.client.Ping(ctx).Err()
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
