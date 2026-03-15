package repository

import (
	"context"
	"fmt"
	"strconv"
	// "tracker_service/internal/models"
	"tracker_service/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DBs struct {
	PDB *pgxpool.Pool
	RDB *redis.Client
}

func NewDBs(pool *pgxpool.Pool, client *redis.Client) *DBs {
	return &DBs{PDB: pool, RDB: client}
}

func (d *DBs) StartWorkout(ctx context.Context, userId int, planId int) (int, error) {
	var trackerId int
	err := d.PDB.QueryRow(ctx, `
		INSERT INTO tracker(user_id, plan_id, started_at)
		VALUES($1, $2, NOW())
		RETURNING id	
	`, userId, planId).Scan(&trackerId)
	if err != nil {
		return trackerId, fmt.Errorf("error starting an empty workout : %w\n", err)
	}

	return trackerId, nil
}
func (d *DBs) RevertStartWorkout(ctx context.Context, trackerId int) error {

	_, err := d.PDB.Exec(ctx, `
		DELETE FROM TRACKER 
		WHERE ID = $1
	`, trackerId)
	if err != nil {
		return fmt.Errorf("error reverting start workout : %w\n", err)
	}

	return nil
}
// func (d *DBs) EndWorkout(ctx context.Context, trackerId int, w models.Tracker) error {
// 	trnx, err := d.PDB.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("error creating transaction : %w\n", err)
// 	}

// 	defer trnx.Rollback(ctx)

// 	for _, allExercises := range w.Workout {

// 		for _, exercise := range allExercises.Tracker {
// 			currentSet := 1
// 			_, err := trnx.Exec(ctx, `
// 				INSERT INTO workout(tracker_id, exercise_id, set_number, weight, reps)
// 				VALUES($1, $2, $3, $4, $5) 
// 			`, trackerId, allExercises.ExerciseId, currentSet, exercise.Weight, exercise.Reps)
// 			if err != nil {
// 				return fmt.Errorf("error inserting workout data for exercise_id %v : %w\n", exercise, err)
// 			}
// 			currentSet = currentSet + 1
// 		}
// 	}
// 	err = trnx.Commit(ctx)
// 	if err != nil {
// 		return fmt.Errorf("error commiting the transaction : %w\n", err)
// 	}

// 	return nil
// }
func (d *DBs) SetTrackerId(ctx context.Context, userId int, trackerId int) error {
	key := fmt.Sprintf("user:%v:tracker_id", userId)
	err := d.RDB.Set(ctx, key, trackerId, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting tracker id : %w", err)
	}
	
	return nil
}
func (d *DBs) DelTrackerId(ctx context.Context, userId int) error {
	key := fmt.Sprintf("user:%v:tracker_id", userId)
	err := d.RDB.Del(ctx, key).Err()
	if err != nil{
		return fmt.Errorf("error in deleting the tracker Id of user with id  %v : %w",userId, err)
	}
	return nil
}
func (d *DBs) GetTrackerId(ctx context.Context, userId int) (int, error) {
	var id int
	key := fmt.Sprintf("user:%v:tracker_id", userId)
	idStr, err := d.RDB.Get(ctx, key).Result()
	if err != nil{
		return id, fmt.Errorf("error in getting the tracker Id of the user with id %v : %w",userId, err)
	}

	id, err = strconv.Atoi(idStr)
	if err != nil{
		return id, fmt.Errorf("error in converting the tracker id from string to int : %w", err)
	}

	return id, nil

}

func (d *DBs) EndWorkout(ctx context.Context, trackerId int, data *user.Tracker) error {

	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error initializing a transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

	// planId := data.PlanId

	for _, allExercises := range data.Workout {
		for _, RepsWeights := range allExercises.RepsWeight {
			currentSet := 1
			_, err := trnx.Exec(ctx, `
				INSERT INTO workout(tracker_id, exercise_id, set_number, weight, reps)
				VALUES($1, $2, $3, $4, $5)			
			`, trackerId, allExercises.ExerciseId, currentSet, RepsWeights.Weight, RepsWeights.Reps)

			if err != nil{
				return fmt.Errorf("error in inserting data into tracker : %w", err)
			}
			currentSet = currentSet + 1
		}
	}


	_, err = trnx.Exec(ctx, `
		UPDATE tracker
		SET ended_at = NOW()
		WHERE id = $1	
	`, trackerId)
	if err != nil{
		return fmt.Errorf("error updating the ended time in tracker : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil{
		return fmt.Errorf("error commiting the transaction : %w", err)
	}

	return nil
}
