package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/sakamoto-max/wt_2_pkg/myerrs"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

func (r *dBs) SetUserEmptyPlanIdR(ctx context.Context, userId string, planId string) error {
	key := fmt.Sprintf("user_id:%v:empty_plan_id", userId)

	err := r.rDB.Set(ctx, key, planId, 0).Err()
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("failed to set empty plan Id in redis : %w", err))
	}

	return nil
}
func (r *dBs) GetUserEmptyPlanIdR(ctx context.Context, userId string) (string, error) {
	key := fmt.Sprintf("user_id:%v:empty_plan_id", userId)

	var planId string

	err := r.rDB.Get(ctx, key).Scan(&planId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", myerrors.InternalServerErrMaker(fmt.Errorf("failed to get empty planId from redis : %w", err))
	}

	return planId, nil
}
func (r *dBs) DelUserEmptyPlanIdR(ctx context.Context, userId string) error {
	key := fmt.Sprintf("user_id:%v:empty_plan_id", userId)

	err := r.rDB.Del(ctx, key).Err()
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("failed to delete the empty planId from redis : %w", err))
	}

	return nil
}

func (r *dBs) SetUserPlanId(ctx context.Context, userId string, planName string, planId string) error {
	key := fmt.Sprintf("user_id:%v:plan_name:%v:id", userId, planName)

	err := r.rDB.Set(ctx, key, planId, 0).Err()
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error setting user plan id for planName %v : %w", planName, err))
	}

	return nil
}

func (r *dBs) GetUserPlanId(ctx context.Context, userId string, planName string) (string, error) {
	key := fmt.Sprintf("user_id:%v:plan_name:%v:id", userId, planName)

	planId, err := r.rDB.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", myerrs.InternalServerErrMaker(fmt.Errorf("error getting planId for user : %w", err))
	}

	return planId, nil

}
func (r *dBs) DelUserPlanId(ctx context.Context, userId string, planName string) error {
	key := fmt.Sprintf("user_id:%v:plan_name:%v:id", userId, planName)

	err := r.rDB.Del(ctx, key).Err()
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error deleting planId : %w", err))
	}

	return nil

}


func (r *dBs) SetUserPlan(ctx context.Context, userId string, planName string, planId string, exerciseIds *[]string) error {

	Key := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)

	pipe := r.rDB.Pipeline()

	for i, exerId := range *exerciseIds {
		exerKey := fmt.Sprintf("exer_%v", i)
		pipe.HSet(ctx, Key, exerKey, exerId)
	}

	pipe.HSet(ctx, Key, "plan_id", planId)
	pipe.HSet(ctx, Key, "number_of_exercises", len(*exerciseIds))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error setting the plan : %w", err))
	}

	return nil
}

func (r *dBs) GetUserPlan(ctx context.Context, userId string, planName string) (string, *[]string, error) {
	Key := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)
	
	res, err := r.rDB.HGetAll(ctx, Key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil, nil
		}
		return "", nil, myerrs.InternalServerErrMaker(fmt.Errorf("error getting user plan from redis : %w", err))
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

func (r *dBs) DelUserPlan(ctx context.Context, userId string, planName string) error {
	Key := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)

	err := r.rDB.Del(ctx, Key).Err()
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error deleting user plan : %w", err))
	}

	return nil
}
