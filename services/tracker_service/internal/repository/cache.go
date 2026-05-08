package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"tracker_service/internal/models"

	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

func (d *dBs) SetTrackerId(ctx context.Context, userId string, trackerId string) error {
	keyforTrackId := fmt.Sprintf("user:%v:tracker_id", userId)

	if err := d.rDB.Set(ctx, keyforTrackId, trackerId, 0).Err(); err != nil {
		return fmt.Errorf("error setting the tracker id : %w", err)
	}

	return nil

}

func (d *dBs) GetTrackerId(ctx context.Context, userId string) (string, error) {
	var id string
	key := fmt.Sprintf("user:%v:tracker_id", userId)
	id, err := d.rDB.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return id, nil
		}
		return id, myerrors.InternalServerErrMaker(fmt.Errorf("error in getting the tracker Id of the user with id %v : %w", userId, err))
	}

	return id, nil

}

func (d *dBs) DelTrackerId(ctx context.Context, userId string) error {
	key := fmt.Sprintf("user:%v:tracker_id", userId)

	if err := d.rDB.Del(ctx, key).Err(); err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting tracker id in redis : %w", err))
	}

	return nil
}

func (d *dBs) SetExerciseNameById(ctx context.Context, exerciseId string, exerciseName string) error {
	key := fmt.Sprintf("exercise_id:%v:name", exerciseId)

	err := d.rDB.Set(ctx, key, exerciseName, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting exercise name : %w", err)
	}

	return nil
}
func (d *dBs) GetExerciseNameById(ctx context.Context, exerciseId string) (string, error) {
	key := fmt.Sprintf("exercise_id:%v:name", exerciseId)

	var exerciseName string
	err := d.rDB.Get(ctx, key).Scan(&exerciseName)
	if err != nil {
		return exerciseName, err
	}

	return exerciseId, nil
}

func (d *dBs) SetUserCurrentPlanName(ctx context.Context, userId string, planName string) error {
	key := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)

	err := d.rDB.Set(ctx, key, planName, 0).Err()

	if err != nil {
		return fmt.Errorf("error setting user current plan : %w", err)
	}

	return nil
}
func (d *dBs) GetUserCurrentPlanName(ctx context.Context, userId string) (string, error) {
	key := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)

	var planName string

	err := d.rDB.Get(ctx, key).Scan(&planName)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", myerrors.InternalServerErrMaker(fmt.Errorf("error getting user current plan : %w", err))
	}

	return planName, nil
}

func (d *dBs) SetPlanWithExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) error {

	key := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)
	noOfExercisesKey := "number_of_exercises"

	pipe := d.rDB.Pipeline()

	for i, exerciseName := range *exerciseNames {
		field := fmt.Sprintf("exer_%v", i)
		pipe.HSet(ctx, key, field, exerciseName)
	}

	pipe.HSet(ctx, key, noOfExercisesKey, len(*exerciseNames))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error in setting plan in redis : %w", err)
	}

	return nil
}
func (d *dBs) GetPlanWithExercises(ctx context.Context, userId string, planName string) (*[]string, error) {

	key := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)
	noOfExercisesKey := "number_of_exercises"

	var NumberOfExercises int
	var numberOfExercisesInString string

	err := d.rDB.HGet(ctx, key, noOfExercisesKey).Scan(&numberOfExercisesInString)
	if err != nil {
		return nil, fmt.Errorf("error getting the number of exercises from redis : %v", err)
	}

	NumberOfExercises, err = strconv.Atoi(numberOfExercisesInString)
	if err != nil {
		return nil, fmt.Errorf("error converting to integer : %w", err)
	}

	var allExercises []string
	var exerciseName string

	for i := 0; i < NumberOfExercises; i++ {
		exerfield := fmt.Sprintf("exer_%v", i)
		err := d.rDB.HGet(ctx, key, exerfield).Scan(&exerciseName)
		if err != nil {
			return nil, fmt.Errorf("unable to get the exercise name : %w", err)
		}

		allExercises = append(allExercises, exerciseName)
	}

	return &allExercises, nil
}
func (d *dBs) SetUserWorkingOutWithPlan(ctx context.Context, userId string, value bool) error {

	key := fmt.Sprintf("user_id:%v:workout_with_plan", userId)

	err := d.rDB.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting user is working out with a plan : %w", err)
	}

	return nil
}
func (d *dBs) GetUserWorkingOutWithPlan(ctx context.Context, userId string) (bool, error) {
	key := fmt.Sprintf("user_id:%v:workout_with_plan", userId)

	cmd := d.rDB.Get(ctx, key)
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

func (d *dBs) SetConflictLevel(ctx context.Context, userId string, conflictLevel int) error {
	key := fmt.Sprintf("user_id:%v:conflict_level", userId)

	err := d.rDB.Set(ctx, key, conflictLevel, 0).Err()
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error setting conflict status : %w", err))
	}

	return nil
}
func (d *dBs) GetConflictLevel(ctx context.Context, userId string) (int, error) {
	key := fmt.Sprintf("user_id:%v:conflict_level", userId)

	var conflictLevel int

	err := d.rDB.Get(ctx, key).Scan(&conflictLevel)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		return 0, myerrors.InternalServerErrMaker(fmt.Errorf("error getting conflict status : %w", err))
	}

	return conflictLevel, nil
}

func (d *dBs) SetUserTrackerData(ctx context.Context, userId string, data *models.Tracker) error {

	pipe := d.rDB.Pipeline()

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
func (d *dBs) GetUserTrackerData(ctx context.Context, userId string) (*models.Tracker, error) {

	mainKey := fmt.Sprintf("user_id:%v:tracker_data", userId)

	numberOfExercisesKey := "number_of_exercises"

	var numberOfExercisesInString string
	err := d.rDB.HGet(ctx, mainKey, numberOfExercisesKey).Scan(&numberOfExercisesInString)
	if err != nil {
		return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting number of exercises from redis : %w", err))
	}

	numberOfExercises, err := strconv.Atoi(numberOfExercisesInString)
	if err != nil {
		return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error converting from number of exercises string to int : %w", err))
	}

	data := models.Tracker{}

	for idForExercise := range numberOfExercises {
		exerciseNameKey := fmt.Sprintf("exercise:%v:name", idForExercise)

		var exerciseName string

		err := d.rDB.HGet(ctx, mainKey, exerciseNameKey).Scan(&exerciseName)
		if err != nil {
			return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting exercise name : %w", err))
		}

		exerciseIdKey := fmt.Sprintf("exercise:%v:id", idForExercise)

		var exerciseId string

		err = d.rDB.HGet(ctx, mainKey, exerciseIdKey).Scan(&exerciseId)
		if err != nil {
			return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting exercise id : %w", err))
		}

		var numberOfSetsString string

		numberOfSetskey := fmt.Sprintf("exercise:%v:number_of_sets", idForExercise)

		err = d.rDB.HGet(ctx, mainKey, numberOfSetskey).Scan(&numberOfSetsString)
		if err != nil {
			return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting exercise id : %w", err))
		}

		numberOfSets, err := strconv.Atoi(numberOfSetsString)
		if err != nil {
			return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error converting number of sets from string to int : %w", err))
		}

		repsAndWeight := []models.Reps{}

		for idForSet := range numberOfSets {

			// repsKey := fmt.Sprintf("exercise:%v:set:%v:reps", exerciseId, idForSet)
			repsKey := fmt.Sprintf("exercise:%v:set:%v:reps", idForExercise, idForSet)

			var repsString string

			err := d.rDB.HGet(ctx, mainKey, repsKey).Scan(&repsString)
			if err != nil {
				return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting reps for exercise %v and set %v : %w", idForExercise, idForSet, err))
			}

			reps, err := strconv.Atoi(repsString)
			if err != nil {
				return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error converting reps from string to int : %w", err))
			}

			weightKey := fmt.Sprintf("exercise:%v:set:%v:weight", idForExercise, idForSet)

			var weightString string

			err = d.rDB.HGet(ctx, mainKey, weightKey).Scan(&weightString)
			if err != nil {
				return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting weight : %w", err))
			}

			// weight, err := strconv.Atoi(weightString)
			weight, err := strconv.ParseFloat(weightString, 32)
			if err != nil {
				return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error converting weight from string to float : %w", err))
			}

			rW := models.Reps{
				Reps:   reps,
				Weight: float32(weight),
			}

			repsAndWeight = append(repsAndWeight, rW)
		}

		w := models.Workout{
			ExerciseId:   exerciseId,
			ExerciseName: exerciseName,
			RepsWeight:   repsAndWeight,
		}

		data.Workout = append(data.Workout, w)
	}

	return &data, nil
}

func (d *dBs) SetUserNewExercises(ctx context.Context, userId string, exerciseNames *[]string) error {

	key := fmt.Sprintf("user_id:%v:new_exercises", userId)
	noOfExercisesKey := "number_of_exercises"

	pipe := d.rDB.Pipeline()

	pipe.HSet(ctx, key, noOfExercisesKey, len(*exerciseNames))

	for i, exerciseName := range *exerciseNames {
		exerciseKey := fmt.Sprintf("exercise_%v", i)
		pipe.HSet(ctx, key, exerciseKey, exerciseName)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
func (d *dBs) GetUserNewExercises(ctx context.Context, userId string) (*[]string, error) {

	key := fmt.Sprintf("user_id:%v:new_exercises", userId)
	noOfExercisesKey := "number_of_exercises"

	var noOfExericisesString string

	err := d.rDB.HGet(ctx, key, noOfExercisesKey).Scan(&noOfExericisesString)
	if err != nil {
		return nil, fmt.Errorf("error getting number of exercises : %w", err)
	}

	noOfExerises, err := strconv.Atoi(noOfExericisesString)
	if err != nil {
		return nil, fmt.Errorf("error converting no of exercises from string to integer : %w", err)
	}

	var allExercises []string

	for i := range noOfExerises {

		var oneExercise string

		exerciseKey := fmt.Sprintf("exercise_%v", i)
		err := d.rDB.HGet(ctx, key, exerciseKey).Scan(&oneExercise)
		if err != nil {
			return nil, fmt.Errorf("error getting the exercise name : %w", err)
		}

		allExercises = append(allExercises, oneExercise)

	}

	return &allExercises, nil
}