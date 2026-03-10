package repository

import (
	"context"
	"fmt"
	"tracker_service/internal/models"

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
		INSERT INTO TRACKER(user_id, plan_id, STARTED_AT)
		VALUES($1, &2, NOW())
		RETURNING id	
	`,userId, planId).Scan(&trackerId)
	if err != nil {
		return trackerId, fmt.Errorf("error starting an empty workout : %w\n", err)
	}

	return trackerId, nil
}
func (d *DBs) RevertStartWorkout(ctx context.Context, trackerId int) (error) {

	_, err := d.PDB.Exec(ctx, `
		DELETE FROM TRACKER 
		WHERE ID = $1
	`, trackerId)
	if err != nil {
		return fmt.Errorf("error reverting start workout : %w\n", err)
	}

	return nil
}

func (d *DBs) EndWorkout(ctx context.Context, trackerId int, w models.Tracker) error {
	trnx, err := d.PDB.Begin(ctx)
	if err != nil{
		return fmt.Errorf("error creating transaction : %w\n", err)
	}

	defer trnx.Rollback(ctx)
	
	for _, allExercises := range w.Workout{

		for _, exercise := range allExercises.Tracker {
			currentSet := 1
			_, err := trnx.Exec(ctx, `
				INSERT INTO workout(tracker_id, exercise_id, set_number, weight, reps)
				VALUES($1, $2, $3, $4, $5) 
			`, trackerId, allExercises.ExerciseId, currentSet, exercise.Weight, exercise.Reps)
			if err != nil{
				return fmt.Errorf("error inserting workout data for exercise_id %v : %w\n", exercise, err)
			}

			currentSet = currentSet + 1
		}

	}

	err = trnx.Commit(ctx)
	if err != nil{
		return fmt.Errorf("error commiting the transaction : %w\n", err)
	}

	return nil
}

func (d *DBs) SetTrackerId(ctx context.Context, userId int, trackerId int) (error) {
	key := fmt.Sprintf("user:%v:tracker_id", userId)
	err := d.RDB.Set(ctx, key, trackerId, 0).Err()
	if err != nil{
		return fmt.Errorf("error setting tracker id : %w", err)
	}

	return nil
}