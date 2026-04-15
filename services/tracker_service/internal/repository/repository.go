package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"tracker_service/internal/models"
	"wt/pkg/enum"
	myerrors "wt/pkg/my_errors"

	"wt/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func (r *DBs) GetPostgresRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.pDB.Ping(ctx)
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
func (r *DBs) GetRedisRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.rDB.Ping(ctx).Err()
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}

func (d *DBs) StartWorkout(ctx context.Context, userId string, planId string) (string, error) {
	var trackerId string
	err := d.pDB.QueryRow(ctx, `
		INSERT INTO tracker(user_id, plan_id, started_at)
		VALUES($1, $2, NOW())
		RETURNING id	
	`, userId, planId).Scan(&trackerId)
	if err != nil {
		return trackerId, myerrors.InternalServerErrMaker(fmt.Errorf("error starting an empty workout : %w\n", err))
	}

	return trackerId, nil
}
func (d *DBs) DeleteTrackerIdInPG(ctx context.Context, trackerId string) error {

	query := `
		DELETE FROM tracker
		WHERE ID = @id	
	`
	_, err := d.pDB.Exec(ctx, query, pgx.NamedArgs{
		"id": trackerId,
	})

	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting the tracker Id : %w", err))
	}

	return nil

}
func (d *DBs) RevertStartWorkout(ctx context.Context, trackerId string) error {

	_, err := d.pDB.Exec(ctx, `
		DELETE FROM TRACKER 
		WHERE ID = $1
	`, trackerId)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error reverting start workout : %w\n", err))
	}

	return nil
}

func (d *DBs) SetTrackerId(ctx context.Context, userId string, trackerId string) error {
	keyforTrackId := fmt.Sprintf("user:%v:tracker_id", userId)

	if err := d.rDB.Set(ctx, keyforTrackId, trackerId, 0).Err(); err != nil {
		return fmt.Errorf("error setting the tracker id : %w", err)
	}

	return nil

}

func (d *DBs) GetTrackerId(ctx context.Context, userId string) (string, error) {
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

func (d *DBs) DelTrackerId(ctx context.Context, userId string) error {
	key := fmt.Sprintf("user:%v:tracker_id", userId)

	if err := d.rDB.Del(ctx, key).Err(); err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting tracker id in redis : %w", err))
	}

	return nil
}

func (d *DBs) EndWorkout(ctx context.Context, trackerId string, data *models.Tracker) error {

	query := `
		INSERT INTO 
			workout(tracker_id, exercise_id, set_number, weight, reps)
		VALUES
			(@tracker_id, @exercise_id, @set_number, @weight, @reps)			
	`

	trnx, err := d.pDB.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error initializing a transaction : %w", err))
	}

	defer trnx.Rollback(ctx)

	for _, dataForEachExercise := range data.Workout {
		exerciseId := dataForEachExercise.ExerciseId
		for _, repsPlusWeight := range dataForEachExercise.RepsWeight {

			currentSet := 1

			_, err := trnx.Exec(ctx, query, pgx.NamedArgs{
				"tracker_id":  trackerId,
				"exercise_id": exerciseId,
				"set_number":  currentSet,
				"weight":      repsPlusWeight.Weight,
				"reps":        repsPlusWeight.Reps,
			})

			if err != nil {
				return myerrors.InternalServerErrMaker(fmt.Errorf("failed to upload data into db : %w", err))
			}

			currentSet = currentSet + 1
		}

	}

	query = `
		UPDATE
			tracker
		SET
			ended_at = NOW()
		WHERE
			id = @tracker_id
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"tracker_id": trackerId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error updating the ended time in tracker : %w", err))
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error commiting the transaction : %w", err))
	}

	return nil
}
func (d *DBs) EndWorkoutWithOutbox(ctx context.Context, userId string, trackerId string, data *models.Tracker, planName string, newExerciseNames *[]string) error {

	query := `
		INSERT INTO 
			workout(tracker_id, exercise_id, set_number, weight, reps)
		VALUES
			(@tracker_id, @exercise_id, @set_number, @weight, @reps)			
	`

	trnx, err := d.pDB.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error initializing a transaction : %w", err))
	}

	defer trnx.Rollback(ctx)

	for _, dataForEachExercise := range data.Workout {
		exerciseId := dataForEachExercise.ExerciseId
		for _, repsPlusWeight := range dataForEachExercise.RepsWeight {

			currentSet := 1

			_, err := trnx.Exec(ctx, query, pgx.NamedArgs{
				"tracker_id":  trackerId,
				"exercise_id": exerciseId,
				"set_number":  currentSet,
				"weight":      repsPlusWeight.Weight,
				"reps":        repsPlusWeight.Reps,
			})

			if err != nil {
				return myerrors.InternalServerErrMaker(fmt.Errorf("failed to upload data into db : %w", err))
			}

			currentSet = currentSet + 1
		}

	}

	query = `
		UPDATE
			tracker
		SET
			ended_at = NOW()
		WHERE
			id = @tracker_id
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"tracker_id": trackerId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error updating the ended time in tracker : %w", err))
	}

	query = `
		INSERT INTO outbox (target_service,	task, payload)
		VALUES (
			@target_service,
			@task,
			@payload
		)
	`

	payload := models.UpdatePlanPayLoad{
		UserId:        userId,
		PlanName:      planName,
		ExerciseNames: newExerciseNames,
	}

	jsonData, err := utils.MakeJSONV2(payload)
	if err != nil {
		return err
	}

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
		"target_service": enum.PlanService,
		"task":           enum.UpdatePlan,
		"payload":        jsonData,
	})

	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error inserting data in outbox : %w", err))
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error commiting the transaction : %w", err))
	}

	return nil
}

func (d *DBs) SetExerciseNameById(ctx context.Context, exerciseId string, exerciseName string) error {
	key := fmt.Sprintf("exercise_id:%v:name", exerciseId)

	err := d.rDB.Set(ctx, key, exerciseName, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting exercise name : %w", err)
	}

	return nil
}
func (d *DBs) GetExerciseNameById(ctx context.Context, exerciseId string) (string, error) {
	key := fmt.Sprintf("exercise_id:%v:name", exerciseId)

	var exerciseName string
	err := d.rDB.Get(ctx, key).Scan(&exerciseName)
	if err != nil {
		return exerciseName, err
	}

	return exerciseId, nil
}

func (d *DBs) SetUserCurrentPlanName(ctx context.Context, userId string, planName string) error {
	key := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)

	err := d.rDB.Set(ctx, key, planName, 0).Err()

	if err != nil {
		return fmt.Errorf("error setting user current plan : %w", err)
	}

	return nil
}
func (d *DBs) GetUserCurrentPlanName(ctx context.Context, userId string) (string, error) {
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

func (d *DBs) SetPlanWithExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) error {

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
func (d *DBs) GetPlanWithExercises(ctx context.Context, userId string, planName string) (*[]string, error) {

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

func (d *DBs) SetUserWorkingOutWithPlan(ctx context.Context, userId string, value bool) error {

	key := fmt.Sprintf("user_id:%v:workout_with_plan", userId)

	err := d.rDB.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting user is working out with a plan : %w", err)
	}

	return nil
}
func (d *DBs) GetUserWorkingOutWithPlan(ctx context.Context, userId string) (bool, error) {
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

func (d *DBs) SetConflictLevel(ctx context.Context, userId string, conflictLevel int) error {
	key := fmt.Sprintf("user_id:%v:conflict_level", userId)

	err := d.rDB.Set(ctx, key, conflictLevel, 0).Err()
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error setting conflict status : %w", err))
	}

	return nil
}
func (d *DBs) GetConflictLevel(ctx context.Context, userId string) (int, error) {
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

func (d *DBs) SetUserTrackerData(ctx context.Context, userId string, data *models.Tracker) error {

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
func (d *DBs) GetUserTrackerData(ctx context.Context, userId string) (*models.Tracker, error) {

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

			weight, err := strconv.Atoi(weightString)
			if err != nil {
				return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error converting weight from string to int : %w", err))
			}

			rW := models.Reps{
				Reps:   reps,
				Weight: weight,
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

func (d *DBs) SetUserNewExercises(ctx context.Context, userId string, exerciseNames *[]string) error {

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
func (d *DBs) GetUserNewExercises(ctx context.Context, userId string) (*[]string, error) {

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

func (d *DBs) DelAllUserData(ctx context.Context, userId string, planName string) error {

	trackerIdKey := fmt.Sprintf("user:%v:tracker_id", userId)
	ongoingWorkoutKey := fmt.Sprintf("user_id:%v:workout_ongoing", userId)
	planWithExercisesKey := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)
	currentPlanKey := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)
	userTrackerDataKey := fmt.Sprintf("user_id:%v:tracker_data", userId)
	newExercisesKey := fmt.Sprintf("user_id:%v:new_exercises", userId)
	conflictKey := fmt.Sprintf("user_id:%v:conflict_level", userId)

	pipe := d.rDB.Pipeline()

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

// func (d *DBs) DelAllUserDataEmptyPlan(ctx context.Context, userId string, planName string) error {

// 	trackerIdKey := fmt.Sprintf("user:%v:tracker_id", userId)
// 	ongoingWorkoutKey := fmt.Sprintf("user_id:%v:workout_ongoing", userId)
// 	// planWithExercisesKey := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)
// 	// currentPlanKey := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)
// 	// userTrackerDataKey := fmt.Sprintf("user_id:%v:tracker_data", userId)
// 	// newExercisesKey := fmt.Sprintf("user_id:%v:new_exercises", userId)
// 	// conflictKey := fmt.Sprintf("user_id:%v:conflict_level", userId)

// 	pipe := d.rDB.Pipeline()

// 	pipe.Del(ctx, trackerIdKey)
// 	pipe.Del(ctx, ongoingWorkoutKey)
// 	pipe.Del(ctx, planWithExercisesKey)
// 	pipe.Del(ctx, userTrackerDataKey)
// 	pipe.Del(ctx, currentPlanKey)
// 	pipe.Del(ctx, newExercisesKey)
// 	pipe.Del(ctx, conflictKey)

// 	_, err := pipe.Exec(ctx)
// 	if err != nil{
// 		return fmt.Errorf("error while deleting userData from redis : %w", err)
// 	}

// 	return nil

// }

func (d *DBs) CancelWorkout(ctx context.Context, userId int) {
	// delete all user data in redis
	// delete the trackerId in postgres
}
