package repository

import (
	"context"
	"errors"
	"exercise_service/internal/models"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repo struct {
	PDB *pgxpool.Pool
	RDB *redis.Client
}

func NewRepo(p *pgxpool.Pool, r *redis.Client) *Repo {
	return &Repo{
		PDB: p,
		RDB: r,
	}
}


func (r *Repo) GetExerciseByName(ctx context.Context, exerciseName string) (*models.Exercise, error) {

	var exercise models.Exercise

	err := r.PDB.QueryRow(ctx, `
		SELECT EXERCISES.ID, EXERCISES.NAME, REST_TIME_IN_SECONDS, 
			BODY_PARTS.NAME, EQUIPMENT.NAME, CREATED_AT FROM EXERCISES
		INNER JOIN 
			BODY_PARTS
		ON 	
			EXERCISES.BODY_PART_ID = BODY_PARTS.ID
		INNER JOIN 
			EQUIPMENT
		ON 
			EXERCISES.EQUIPMENT_ID = EQUIPMENT.ID
		WHERE 
			EXERCISES.name = $1
	`, exerciseName).
	Scan(&exercise.Id, &exercise.Name, &exercise.RestTime, &exercise.BodyPart, 
		&exercise.Equipment, &exercise.CreatedAt) 
	if err != nil{
		return &exercise, err
	}

	return &exercise, nil
}
func (r *Repo) GetAllExercises(ctx context.Context) (*[]models.Exercise, error) {

	var allExercises []models.Exercise
	// var exercise models.Exercise

	rows, err := r.PDB.Query(ctx, `
		SELECT EXERCISES.ID, EXERCISES.NAME, REST_TIME_IN_SECONDS, 
			BODY_PARTS.NAME, EQUIPMENT.NAME, CREATED_AT FROM EXERCISES
		INNER JOIN 
			BODY_PARTS
		ON 	
			EXERCISES.BODY_PART_ID = BODY_PARTS.ID
		INNER JOIN 
			EQUIPMENT
		ON 
			EXERCISES.EQUIPMENT_ID = EQUIPMENT.ID
	`)
	if err != nil{
		return &allExercises, fmt.Errorf("error getting all exercises from DB : %w\n", err)
	}

	var id int
	var exerciseName string
	var restTime int
	var bodyPart string
	var equipmentName string
	var createdAt time.Time
	for rows.Next() {
		err := rows.Scan(&id, &exerciseName, &restTime, &bodyPart, &equipmentName, &createdAt)
		if err != nil{
			return &allExercises, fmt.Errorf("error scanning the rows : %w\n", err)
		}

		exercise := models.Exercise{
			Id: id, 
			Name: exerciseName, 
			RestTime: restTime, 
			BodyPart: bodyPart, 
			Equipment:  equipmentName, 
			CreatedAt: createdAt, 
		}

		allExercises = append(allExercises, exercise)
	}

	return &allExercises, nil
}


func (r *Repo) DeleteExecise(ctx context.Context, exerciseName string) error {
	_, err := r.PDB.Exec(ctx, `
		DELETE FROM EXERCISES
		WHERE name = $1
	`, exerciseName)

	if err != nil{
		return fmt.Errorf("error deleting the exercise %v : %w", exerciseName, err)
	}

	return nil
}

func (r *Repo) CreateExercise(ctx context.Context, exercise *models.Exercise) error {

	var bodyPartId int
	var equipmentId int
	trnx, err := r.PDB.Begin(ctx)
	if err != nil{
		return fmt.Errorf("error creating a transaction : %w\n", err)
	}

	defer trnx.Rollback(ctx)

	err = trnx.QueryRow(ctx, `
		SELECT ID FROM BODY_PARTS
		WHERE NAME = $1	
	`, exercise.BodyPart).Scan(&bodyPartId)
	if err != nil {
		return fmt.Errorf("error getting id of body_part : %w\n", err)
	}

	err = trnx.QueryRow(ctx, `
		SELECT ID FROM EQUIPMENT
		WHERE NAME = $1	
	`, exercise.Equipment).Scan(&equipmentId)
	if err != nil {
		return fmt.Errorf("error getting id of equipment : %w\n", err)
	}

	_, err = trnx.Exec(ctx, `
		INSERT INTO EXERCISES(name, body_part_id, rest_time_in_seconds, equipment_id, created_at)
		VALUES($1, $2, $3, $4, NOW())	
	`, exercise.Name, bodyPartId, exercise.RestTime, equipmentId)
	if err != nil{
		return fmt.Errorf("error inserting exercise %v : %w\n", exercise.Name, err)
	}

	err = trnx.Commit(ctx)
	if err != nil{
		return fmt.Errorf("error commiting the transaction : %w\n", err)
	}

	return nil
}

func (r *Repo) ExerciseExistsReturnId(ctx context.Context, exerciseName string) (bool, int32, error) {

	var id int32
	err := r.PDB.QueryRow(ctx, `
		SELECT id FROM exercises
		WHERE name = $1
	`, exerciseName).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return false, id, nil
		}

		return false, id, fmt.Errorf("error checking if the exericse exists : %v", err)
	}

	return true, id, nil
}