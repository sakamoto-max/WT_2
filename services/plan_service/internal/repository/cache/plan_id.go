package cache

import (
	"context"
	"errors"
	"fmt"
	"plan_service/internal/domain"

	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type planIdCache struct {
	client *redis.Client
}

func (c *planIdCache) SetUserPlanId(ctx context.Context, payload domain.GetPlan, planId string) {
	key := fmt.Sprintf("user_id:%v:plan_name:%v:id", payload.UserId, payload.PlanName)

	c.client.Set(ctx, key, planId, 0)
	// if err != nil {
	// 	// return myerrors.InternalServerErrMaker(fmt.Errorf("error setting user plan id for planName %v : %w", planName, err))
	// }
}
func (c *planIdCache) GetUserPlanId(ctx context.Context, payload domain.GetPlan) (string, error) {
	key := fmt.Sprintf("user_id:%v:plan_name:%v:id", payload.UserId, payload.PlanName)

	planId, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", myerrors.InternalServerErrMaker(fmt.Errorf("error getting planId for user : %w", err))
	}

	return planId, nil

}
func (c *planIdCache) DelUserPlanId(ctx context.Context, payload domain.GetPlan) {
	key := fmt.Sprintf("user_id:%v:plan_name:%v:id", payload.UserId, payload.PlanName)

	c.client.Del(ctx, key)
}
