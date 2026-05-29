package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type emptyPlanCache struct {
	client *redis.Client
}

func (c *emptyPlanCache) SetUserEmptyPlanId(ctx context.Context, userId string, planId string) {
	key := fmt.Sprintf("user_id:%v:empty_plan_id", userId)

	c.client.Set(ctx, key, planId, 0).Err()

}
func (c *emptyPlanCache) GetUserEmptyPlanId(ctx context.Context, userId string) (string, error) {
	key := fmt.Sprintf("user_id:%v:empty_plan_id", userId)

	var planId string

	err := c.client.Get(ctx, key).Scan(&planId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", myerrors.InternalServerErrMaker(fmt.Errorf("failed to get empty planId from redis : %w", err))
	}

	return planId, nil
}
func (c *emptyPlanCache) DelUserEmptyPlanId(ctx context.Context, userId string) error {
	key := fmt.Sprintf("user_id:%v:empty_plan_id", userId)

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("failed to delete the empty planId from redis : %w", err))
	}

	return nil
}
