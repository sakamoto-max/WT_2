package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"tracker_service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
)

type endWorkoutRepo struct {
	pg *pgxpool.Pool
}

func (d *endWorkoutRepo) EndWorkout(ctx context.Context, trackerId string, data *domain.Tracker) error {

	query := `
		INSERT INTO 
			workout(tracker_id, exercise_id, set_number, weight, reps)
		VALUES
			(@tracker_id, @exercise_id, @set_number, @weight, @reps)			
	`
	trnx, err := d.pg.Begin(ctx)
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
func (d *endWorkoutRepo) EndWorkoutWithOutbox(ctx context.Context, userId string, trackerId string, data *domain.Tracker, planName string, newExerciseNames *[]string) error {

	query := `
		INSERT INTO 
			workout(tracker_id, exercise_id, set_number, weight, reps)
		VALUES
			(@tracker_id, @exercise_id, @set_number, @weight, @reps)			
	`

	trnx, err := d.pg.Begin(ctx)
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
		INSERT INTO outbox (target_service, created_by, task, payload)
		VALUES (
			@target_service,
			@createdBy,
			@task,
			@payload
		)
	`

	payload := map[string]any{
		enum.QueueFields_USER_ID.String():        userId,
		enum.QueueFields_PLAN_NAME.String():      planName,
		enum.QueueFields_EXERCISE_NAMES.String(): newExerciseNames,
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
		"createdBy" : enum.ServiceName_TRACKER_SERVICE.String(),
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
