package cache

import (
	"context"
	"errors"
	"fmt"
	"plan_service/internal/domain"
	"strconv"

	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type userPlanCache struct {
	client *redis.Client
}

func (c *userPlanCache) SetUserPlan(ctx context.Context, payload domain.Plan) {

	Key := fmt.Sprintf("user_id:%v:plan_name:%v", payload.UserId, payload.PlanName)

	pipe := c.client.Pipeline()

	for i, exerId := range *payload.ExerciseIds {
		exerKey := fmt.Sprintf("exer_%v", i)
		pipe.HSet(ctx, Key, exerKey, exerId)
	}

	pipe.HSet(ctx, Key, "plan_id", payload.PlanId)
	pipe.HSet(ctx, Key, "number_of_exercises", len(*payload.ExerciseIds))

	pipe.Exec(ctx)
}

func (c *userPlanCache) GetUserPlan(ctx context.Context, payload domain.GetPlan) (string, *[]string, error) {
	Key := fmt.Sprintf("user_id:%v:plan_name:%v", payload.UserId, payload.PlanName)

	res, err := c.client.HGetAll(ctx, Key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil, nil
		}
		return "", nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting user plan from redis : %w", err))
	}

	numberOfExerciseStr := res["number_of_exercises"]

	numberOfExercises, _ := strconv.Atoi(numberOfExerciseStr)

	var ExerIds []string

	for i := range numberOfExercises {
		exerKey := fmt.Sprintf("exer_%v", i)

		ExerIds = append(ExerIds, res[exerKey])
	}

	planId := res["plan_id"]

	return planId, &ExerIds, nil
}

func (c *userPlanCache) DelUserPlan(ctx context.Context, payload domain.GetPlan) error {
	Key := fmt.Sprintf("user_id:%v:plan_name:%v", payload.UserId, payload.PlanName)

	cmd := c.client.Del(ctx, Key)

	if cmd.Err() != nil {
		return fmt.Errorf("error deleting the user plan from cache : %w", cmd.Err())
	}

	return nil

}
