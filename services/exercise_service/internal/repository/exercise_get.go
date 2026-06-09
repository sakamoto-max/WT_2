package repository

import (
	"context"
	"errors"
	"exercise_service/internal/domain"
	"exercise_service/internal/mappings"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type exerciseGetDB struct {
	pg *pgxpool.Pool
}

func (d *exerciseGetDB) GetAllExercises(ctx context.Context, payload mappings.GetAllExercises) (*[]domain.Exercise, error) {

	query := `
		SELECT 
			EXERCISES.ID, 
			EXERCISES.NAME, 
			BODY_PARTS.NAME, 
			EQUIPMENT.NAME,
			created_at,
			updated_at
		FROM 
			EXERCISES
		INNER JOIN 
			BODY_PARTS 
		ON 
			EXERCISES.BODY_PART_ID = BODY_PARTS.ID
		INNER JOIN 
			EQUIPMENT 
		ON 
			EXERCISES.EQUIPMENT_ID = EQUIPMENT.ID
		FULL JOIN (
			SELECT * FROM USER_NULLIFIED
			WHERE USER_ID = @userId
			) AS TABLE_B
		ON 
			EXERCISES.ID = TABLE_B.EXERCISE_ID
		WHERE 
			(CREATED_BY IS NULL OR CREATED_BY = @userId) AND TABLE_B.EXERCISE_ID IS NULL
	`

	if payload.BodyPart != "" {
		query += ` AND BODY_PARTS.NAME = @bodypart`
	}

	if payload.Equipment != "" {
		query += ` AND EQUIPMENT.NAME = @equipment`
	}

	
	rows, err := d.pg.Query(ctx, query, pgx.NamedArgs{
		"userId":    payload.UserId,
		"bodypart":  payload.BodyPart,
		"equipment": payload.Equipment,
	})
	if err != nil {
		
		err := fmt.Errorf("error getting all exercises from DB : %w\n", err)
		
		return nil, myerrs.InternalServerErrMaker(err)
	}
	var allExercises []domain.Exercise

	var id string
	var exerciseName string
	var bodyPart string
	var equipmentName string
	var createdAt time.Time
	var updatedAt time.Time

	for rows.Next() {
		err := rows.Scan(&id, &exerciseName, &bodyPart, &equipmentName, &createdAt, &updatedAt)
		if err != nil {
			err := fmt.Errorf("error scanning the rows : %w\n", err)
			return nil, myerrs.InternalServerErrMaker(err)
		}

		exercise := domain.Exercise{
			Id:        id,
			Name:      exerciseName,
			BodyPart:  bodyPart,
			Equipment: equipmentName,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		allExercises = append(allExercises, exercise)
	}



	if len(allExercises) == 0 {
		
		stmt := "exercises"
		
		if payload.BodyPart != "" {
			stmt += fmt.Sprintf(" with body part %s,", payload.BodyPart)
		}
		
		if payload.Equipment != "" {
			stmt += fmt.Sprintf(" with equipment name %s", payload.Equipment)
		}
		
		// stmt += " not found"
		
		return nil, myerrs.ResourceNotFoundErrMaker(stmt)

	}

	return &allExercises, nil
}
func (d *exerciseGetDB) GetExerciseByName(ctx context.Context, payload mappings.GetExerciseByName) (*domain.Exercise, error) {

	var exercise domain.Exercise
	query := `
			SELECT 
			EXERCISES.ID, 
			EXERCISES.NAME, 
			BODY_PARTS.NAME, 
			EQUIPMENT.NAME ,
			created_at,
			updated_at
		FROM 
			EXERCISES
		INNER JOIN 
			BODY_PARTS 
		ON 
			EXERCISES.BODY_PART_ID = BODY_PARTS.ID
		INNER JOIN 
			EQUIPMENT 
		ON 
			EXERCISES.EQUIPMENT_ID = EQUIPMENT.ID
		FULL JOIN (
			SELECT * FROM USER_NULLIFIED
			WHERE USER_ID = @userId
			) AS TABLE_B
		ON 
			EXERCISES.ID = TABLE_B.EXERCISE_ID
		WHERE 
			Exercises.NAME = @exerciseName and (CREATED_BY IS NULL OR CREATED_BY = @userId) AND TABLE_B.EXERCISE_ID IS NULL;	
	`
	row := d.pg.QueryRow(ctx, query, pgx.NamedArgs{"exerciseName": payload.ExerciseName, "userId": payload.UserId})
	err := row.Scan(&exercise.Id, &exercise.Name, &exercise.BodyPart, &exercise.Equipment, &exercise.CreatedAt, &exercise.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrs.ResourceNotFoundErrMaker("exercise")
		}

		return nil, myerrs.InternalServerErrMaker(fmt.Errorf("error in getting the exercise by name : %v : %w", payload.ExerciseName, err))
	}

	return &exercise, nil
}
func (d *exerciseGetDB) GetExerciseNameByID(ctx context.Context, exerciseId string) (string, error) {

	var exerciseName string

	err := d.pg.QueryRow(ctx, `
		SELECT name from exercises
		where id = @id	
	`, pgx.NamedArgs{"id": exerciseId}).Scan(&exerciseName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", myerrs.ResourceNotFoundErrMaker("exericse")
		}

		err := fmt.Errorf("error getting execise name for id %v : %w", exerciseId, err)
		return "", myerrs.InternalServerErrMaker(err)
	}

	return exerciseName, nil
}
func (d *exerciseGetDB) ExerciseExistsReturnId(ctx context.Context, payload mappings.ExerciseExistsReturnId) (string, error) {

	query := `
		SELECT 
			ID 
		FROM 
			EXERCISES
		LEFT JOIN 
			(SELECT EXERCISE_ID FROM USER_NULLIFIED WHERE USER_ID = @userId) AS TABLE_2
		ON 
			EXERCISES.ID = TABLE_2.EXERCISE_ID
		WHERE 
			NAME = @exerciseName
		AND 
			(CREATED_BY = @userId OR CREATED_BY IS NULL) 
		AND 
			TABLE_2.EXERCISE_ID IS NULL;	
	`
	var id string
	err := d.pg.QueryRow(ctx, query, pgx.NamedArgs{"exerciseName": payload.ExerciseName, "userId": payload.UserId}).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return id, myerrs.ResourceNotFoundErrMaker(payload.ExerciseName)
		}
		return id, myerrs.InternalServerErrMaker(fmt.Errorf("error checking if the exericse exists : %v", err))
	}

	return id, nil
}
