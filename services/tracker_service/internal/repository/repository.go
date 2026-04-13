package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tracker_service/internal/models"
	myerrors "wt/pkg/my_errors"

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

func (d *DBs) SetTrackerIdAndOngoingWorkout(ctx context.Context, userId string, trackerId string) error {
	keyforTrackId := fmt.Sprintf("user:%v:tracker_id", userId)
	keyforOngoingWorkout := fmt.Sprintf("user_id:%v:workout_ongoing", userId)

	pipe := d.rDB.Pipeline()

	pipe.Set(ctx, keyforTrackId, trackerId, 0)
	pipe.Set(ctx, keyforOngoingWorkout, true, 0)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error setting the tracker Id and ongoing workout : %w", err))
	}

	return nil
}
func (d *DBs) DelTrackerIdAndOngoingWorkout(ctx context.Context, userId string) error {

	key := fmt.Sprintf("user:%v:tracker_id", userId)
	keyforOngoingWorkout := fmt.Sprintf("user_id:%v:workout_ongoing", userId)

	pipe := d.rDB.Pipeline()

	pipe.Del(ctx, keyforOngoingWorkout)
	pipe.Del(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error in deleting the tracker Id and ongoing workout of user with id  %v : %w", userId, err))
	}
	return nil
}
func (d *DBs) GetTrackerId(ctx context.Context, userId string) (string, error) {
	var id string
	key := fmt.Sprintf("user:%v:tracker_id", userId)
	id, err := d.rDB.Get(ctx, key).Result()
	if err != nil {
		return id, myerrors.InternalServerErrMaker(fmt.Errorf("error in getting the tracker Id of the user with id %v : %w", userId, err))
	}

	return id, nil

}

func (d *DBs) CheckIfWorkoutIsOngoing(ctx context.Context, userId string) (bool, error) {
	keyforOngoingWorkout := fmt.Sprintf("user_id:%v:workout_ongoing", userId)

	res, err := d.rDB.Get(ctx, keyforOngoingWorkout).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, myerrors.InternalServerErrMaker(fmt.Errorf("error in checking if user has ongoing workout : %w", err))
	}

	if res == "0" {
		return false, nil
	}

	return true, nil
}

func (d *DBs) DelOngoingWorkout() {

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

			if err != nil{
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

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"tracker_id" : trackerId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error updating the ended time in tracker : %w", err))
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

func (d *DBs) SetUserCurrentPlanName(ctx context.Context, userId int, planName string) error {
	key := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)

	err := d.rDB.Set(ctx, key, planName, 0).Err()

	if err != nil{
		return fmt.Errorf("error setting user current plan : %w", err)
	}

	return nil
}
func (d *DBs) GetUserCurrentPlanName(ctx context.Context, userId int) (string, error) {
	key := fmt.Sprintf("user_id:%v:current_workout_plan_name", userId)

	var planName string

	err := d.rDB.Get(ctx, key).Scan(&planName)

	if err != nil{
		return "", fmt.Errorf("error getting user current plan : %w", err)
	}

	return planName, nil
}

func (d *DBs) SetPlanWithExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) error {
	// key : user_id:%v:plan_name:%v
	// field : exer_%v
	// val : exercise_name

	key := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)
	noOfExercisesKey := "number_of_exercises"
	
	pipe := d.rDB.Pipeline()
	
	for i, exerciseName := range *exerciseNames{
		field := fmt.Sprintf("exer_%v", i)
		pipe.HSet(ctx, key, field, exerciseName)
	}
	
	pipe.HSet(ctx, key, noOfExercisesKey, len(*exerciseNames))
	
	_, err := pipe.Exec(ctx)
	if err != nil{
		return fmt.Errorf("error in setting plan in redis : %w", err)
	}

	return nil
}

// func (d *DBs) GetPlanWithExercises(ctx context.Context, userId int, planName string) {
// 	key := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)

// 	cmd := d.rDB.HGetAll(ctx, key)
// 	a, _ := cmd.Result()
// 	fmt.Println(a)
// }

func (d *DBs) SetUserWorkingOutWithPlan(ctx context.Context, userId string, value bool) error {

	key := fmt.Sprintf("user_id:%v:workout_with_plan", userId)
	
	err := d.rDB.Set(ctx, key, value, 0).Err()
	if err != nil{
		return fmt.Errorf("error setting user is working out with a plan : %w", err)
	}
	
	return nil
}

func (d *DBs) GetUserWorkingOutWithPlan(ctx context.Context, userId string) (bool, error){
	key := fmt.Sprintf("user_id:%v:workout_with_plan", userId)

	cmd := d.rDB.Get(ctx, key)
	res, err := cmd.Result()

	if err != nil{
		if errors.Is(err, redis.Nil){
			return false, nil
		}
		return false,  fmt.Errorf("error getting user working out with plan : %w", err)
	}

	if res == "false" {
		return false, nil
	}

	return true, nil
	
}

func (d *DBs) DelUserWorkingOutWithPlan(ctx context.Context, userId string) error {
	key := fmt.Sprintf("user_id:%v:workout_with_plan", userId)

	err := d.rDB.Del(ctx, key).Err()
	if err != nil{
		return fmt.Errorf("error deleting user working out with plan : %w", err)
	}

	return nil
}