package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type CurrentPlanCache struct {
	client *redis.Client
}

func (c *CurrentPlanCache) SetUserCurrentPlanName(ctx context.Context, userId string, planName string) error {
	key := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)

	err := c.client.Set(ctx, key, planName, 0).Err()

	if err != nil {
		return fmt.Errorf("error setting user current plan : %w", err)
	}

	return nil
}
func (c *CurrentPlanCache) GetUserCurrentPlanName(ctx context.Context, userId string) (string, error) {
	key := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)

	var planName string

	err := c.client.Get(ctx, key).Scan(&planName)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", myerrors.InternalServerErrMaker(fmt.Errorf("error getting user current plan : %w", err))
	}

	return planName, nil
}
