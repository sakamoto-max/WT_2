package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type userDataCache struct {
	client *redis.Client
}

func (c *userDataCache) DelAllUserData(ctx context.Context, userId string, planName string) error {

	trackerIdKey := fmt.Sprintf("user:%v:tracker_id", userId)
	ongoingWorkoutKey := fmt.Sprintf("user_id:%v:workout_ongoing", userId)
	planWithExercisesKey := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)
	currentPlanKey := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)
	userTrackerDataKey := fmt.Sprintf("user_id:%v:tracker_data", userId)
	newExercisesKey := fmt.Sprintf("user_id:%v:new_exercises", userId)
	conflictKey := fmt.Sprintf("user_id:%v:conflict_level", userId)

	pipe := c.client.Pipeline()

	pipe.Del(ctx, trackerIdKey)
	pipe.Del(ctx, ongoingWorkoutKey)
	pipe.Del(ctx, planWithExercisesKey)
	pipe.Del(ctx, userTrackerDataKey)
	pipe.Del(ctx, currentPlanKey)
	pipe.Del(ctx, newExercisesKey)
	pipe.Del(ctx, conflictKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error while deleting userData from redis : %w", err))
	}

	return nil

}
