package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2-pkg/my_errors"
)

func (r *DBs) SetUserEmptyPlanIdR(ctx context.Context, userId string, planId string) error {
	key := fmt.Sprintf("user_id:%v:empty_plan_id", userId)

	err := r.rDB.Set(ctx, key, planId, 0).Err()
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("failed to set empty plan Id in redis : %w", err))
	}

	return nil
}
func (r *DBs) GetUserEmptyPlanIdR(ctx context.Context, userId string) (string, error) {
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
func (r *DBs) DelUserEmptyPlanIdR(ctx context.Context, userId string) error {
	key := fmt.Sprintf("user_id:%v:empty_plan_id", userId)

	err := r.rDB.Del(ctx, key).Err()
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("failed to delete the empty planId from redis : %w", err))
	}

	return nil
}
