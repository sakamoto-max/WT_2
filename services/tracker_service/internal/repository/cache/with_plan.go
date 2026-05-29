package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type withPlanCache struct {
	client *redis.Client
}

func (c *withPlanCache) SetUserWorkingOutWithPlan(ctx context.Context, userId string, value bool) error {

	key := fmt.Sprintf("user_id:%v:workout_with_plan", userId)

	err := c.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting user is working out with a plan : %w", err)
	}

	return nil
}
func (c *withPlanCache) GetUserWorkingOutWithPlan(ctx context.Context, userId string) (bool, error) {
	key := fmt.Sprintf("user_id:%v:workout_with_plan", userId)

	cmd := c.client.Get(ctx, key)
	res, err := cmd.Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("error getting user working out with plan : %w", err)
	}

	if res == "false" {
		return false, nil
	}

	return true, nil

}
