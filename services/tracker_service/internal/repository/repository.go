package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"tracker_service/internal/models"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/jackc/pgx/v5"
)

func (r *dBs) GetPostgresRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.pDB.Ping(ctx)
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
func (r *dBs) GetRedisRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.rDB.Ping(ctx).Err()
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}

func (d *dBs) StartWorkout(ctx context.Context, userId string, planId string) (string, error) {
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
func (d *dBs) DeleteTrackerIdInPG(ctx context.Context, trackerId string) error {

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
func (d *dBs) RevertStartWorkout(ctx context.Context, trackerId string) error {

	_, err := d.pDB.Exec(ctx, `
		DELETE FROM TRACKER 
		WHERE ID = $1
	`, trackerId)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error reverting start workout : %w\n", err))
	}

	return nil
}

func (d *dBs) EndWorkout(ctx context.Context, trackerId string, data *models.Tracker) error {

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
		for i, repsPlusWeight := range dataForEachExercise.RepsWeight {

			_, err := trnx.Exec(ctx, query, pgx.NamedArgs{
				"tracker_id":  trackerId,
				"exercise_id": exerciseId,
				"set_number":  i + 1,
				"weight":      repsPlusWeight.Weight,
				"reps":        repsPlusWeight.Reps,
			})

			if err != nil {
				return myerrors.InternalServerErrMaker(fmt.Errorf("failed to upload data into db : %w", err))
			}
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
func (d *dBs) EndWorkoutWithOutbox(ctx context.Context, userId string, trackerId string, data *models.Tracker, planName string, newExerciseNames *[]string) error {

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

	dataInBytes, err := json.Marshal(payload)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("failed to marshal the data : %w", err))
	}

	jsonData := string(dataInBytes)

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
		"target_service": enum.ServiceName_PLAN_SERVICE.String(),
		"task":           enum.TaskName_UPDATE_PLAN.String(),
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






func (d *dBs) DelAllUserData(ctx context.Context, userId string, planName string) error {

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

// func (d *dBs) DelAllUserDataEmptyPlan(ctx context.Context, userId string, planName string) error {

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

func (d *dBs) CancelWorkout(ctx context.Context, userId int) {
	// delete all user data in redis
	// delete the trackerId in postgres
}
