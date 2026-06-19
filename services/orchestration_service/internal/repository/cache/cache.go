package cache

import (
	"context"
	"errors"
	"fmt"
	"orchestration_service/internal/types"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		redis: client,
	}
}

func (c *Cache) SetTaskTimeOut(ctx context.Context, data types.Data) error {

	key := fmt.Sprintf("task_id:%v:timeout", data.TaskId)

	err := c.redis.Set(ctx, key, true, data.TimeOut.Duration).Err()
	if err != nil {
		return fmt.Errorf("failed to set task timeout : %w", err)
	}

	return nil
}

func (c *Cache) SkipTask(ctx context.Context, data types.Data) (bool, error) {
	key := fmt.Sprintf("task_id:%v:timeout", data.TaskId)

	var value string

	err := c.redis.Get(ctx, key).Scan(&value)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get the value from redis : %w", err)
	}

	if value == "true" {
		return true, nil
	}

	return false, nil
}
