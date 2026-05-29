package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type exerciseNameCache struct {
	client *redis.Client
}

func (c *exerciseNameCache) SetExerciseNameById(ctx context.Context, exerciseId string, exerciseName string) error {
	key := fmt.Sprintf("exercise_id:%v:name", exerciseId)

	err := c.client.Set(ctx, key, exerciseName, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting exercise name : %w", err)
	}

	return nil
}
func (c *exerciseNameCache) GetExerciseNameById(ctx context.Context, exerciseId string) (string, error) {
	key := fmt.Sprintf("exercise_id:%v:name", exerciseId)

	var exerciseName string
	err := c.client.Get(ctx, key).Scan(&exerciseName)
	if err != nil {
		return exerciseName, err
	}

	return exerciseId, nil
}
