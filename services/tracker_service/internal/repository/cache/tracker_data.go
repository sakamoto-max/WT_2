package cache

import (
	"context"
	"fmt"
	"strconv"
	"tracker_service/internal/domain"
	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type trackerDataCache struct {
	client *redis.Client
}

func (c *trackerDataCache) SetUserTrackerData(ctx context.Context, userId string, data *domain.Tracker) error {

	pipe := c.client.Pipeline()

	mainKey := fmt.Sprintf("user_id:%v:tracker_data", userId)
	numberOfExercises := "number_of_exercises"

	pipe.HSet(ctx, mainKey, numberOfExercises, len(data.Workout))

	for idForExercise, exercise := range data.Workout {
		exerciseNameKey := fmt.Sprintf("exercise:%v:name", idForExercise)
		pipe.HSet(ctx, mainKey, exerciseNameKey, exercise.ExerciseName)

		exerciseIdKey := fmt.Sprintf("exercise:%v:id", idForExercise)
		pipe.HSet(ctx, mainKey, exerciseIdKey, exercise.ExerciseId)

		numberOfSets := fmt.Sprintf("exercise:%v:number_of_sets", idForExercise)
		pipe.HSet(ctx, mainKey, numberOfSets, len(exercise.RepsWeight))

		for idForSet, repsAndWeight := range exercise.RepsWeight {
			repsKey := fmt.Sprintf("exercise:%v:set:%v:reps", idForExercise, idForSet)
			pipe.HSet(ctx, mainKey, repsKey, repsAndWeight.Reps)

			weightKey := fmt.Sprintf("exercise:%v:set:%v:weight", idForExercise, idForSet)
			pipe.HSet(ctx, mainKey, weightKey, repsAndWeight.Weight)
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error setting user tracker data in redis : %w", err))
	}

	return nil
}
func (c *trackerDataCache) GetUserTrackerData(ctx context.Context, userId string) (*domain.Tracker, error) {

	mainKey := fmt.Sprintf("user_id:%v:tracker_data", userId)

	numberOfExercisesKey := "number_of_exercises"

	var numberOfExercisesInString string
	err := c.client.HGet(ctx, mainKey, numberOfExercisesKey).Scan(&numberOfExercisesInString)
	if err != nil {
		return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting number of exercises from redis : %w", err))
	}

	numberOfExercises, err := strconv.Atoi(numberOfExercisesInString)
	if err != nil {
		return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error converting from number of exercises string to int : %w", err))
	}

	data := domain.Tracker{}

	for idForExercise := range numberOfExercises {
		exerciseNameKey := fmt.Sprintf("exercise:%v:name", idForExercise)

		var exerciseName string

		err := c.client.HGet(ctx, mainKey, exerciseNameKey).Scan(&exerciseName)
		if err != nil {
			return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting exercise name : %w", err))
		}

		exerciseIdKey := fmt.Sprintf("exercise:%v:id", idForExercise)

		var exerciseId string

		err = c.client.HGet(ctx, mainKey, exerciseIdKey).Scan(&exerciseId)
		if err != nil {
			return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting exercise id : %w", err))
		}

		var numberOfSetsString string

		numberOfSetskey := fmt.Sprintf("exercise:%v:number_of_sets", idForExercise)

		err = c.client.HGet(ctx, mainKey, numberOfSetskey).Scan(&numberOfSetsString)
		if err != nil {
			return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting exercise id : %w", err))
		}

		numberOfSets, err := strconv.Atoi(numberOfSetsString)
		if err != nil {
			return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error converting number of sets from string to int : %w", err))
		}

		repsAndWeight := []domain.Reps{}

		for idForSet := range numberOfSets {

			// repsKey := fmt.Sprintf("exercise:%v:set:%v:reps", exerciseId, idForSet)
			repsKey := fmt.Sprintf("exercise:%v:set:%v:reps", idForExercise, idForSet)

			var repsString string

			err := c.client.HGet(ctx, mainKey, repsKey).Scan(&repsString)
			if err != nil {
				return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting reps for exercise %v and set %v : %w", idForExercise, idForSet, err))
			}

			reps, err := strconv.Atoi(repsString)
			if err != nil {
				return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error converting reps from string to int : %w", err))
			}

			weightKey := fmt.Sprintf("exercise:%v:set:%v:weight", idForExercise, idForSet)

			var weightString string

			err = c.client.HGet(ctx, mainKey, weightKey).Scan(&weightString)
			if err != nil {
				return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting weight : %w", err))
			}

			// weight, err := strconv.Atoi(weightString)
			weight, err := strconv.ParseFloat(weightString, 32)
			if err != nil {
				return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error converting weight from string to float : %w", err))
			}

			rW := domain.Reps{
				Reps:   reps,
				Weight: float32(weight),
			}

			repsAndWeight = append(repsAndWeight, rW)
		}

		w := domain.Workout{
			ExerciseId:   exerciseId,
			ExerciseName: exerciseName,
			RepsWeight:   repsAndWeight,
		}

		data.Workout = append(data.Workout, w)
	}

	return &data, nil
}
