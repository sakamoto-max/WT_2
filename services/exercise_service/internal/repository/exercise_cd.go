package repository

import (
	"context"
	"errors"
	"exercise_service/mappings"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type exerciseCD struct {
	pg *pgxpool.Pool
}

func (d *exerciseCD) CreateExercise(ctx context.Context, payload mappings.CreateExercise) (string, error) {

	var bodyPartId uuid.UUID
	var equipmentId uuid.UUID
	var Id uuid.UUID

	trnx, err := d.pg.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating a transaction : %w\n", err)
	}

	defer trnx.Rollback(ctx)

	query := `
		SELECT ID FROM BODY_PARTS
		WHERE NAME = @name
	`

	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{"name": payload.BodyPartName}).Scan(&bodyPartId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", myerrs.ResourceNotFoundErrMaker("body_part")
		}

		err := fmt.Errorf("error getting id of body_part : %w\n", err)
		return "", myerrs.InternalServerErrMaker(err)
	}

	query = `
		SELECT ID FROM EQUIPMENT
		WHERE NAME = @name
	`

	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{"name": payload.EquipmentName}).Scan(&equipmentId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", myerrs.ResourceNotFoundErrMaker("equipment")
		}

		err := fmt.Errorf("error getting id of equipment : %w\n", err)
		return "", myerrs.InternalServerErrMaker(err)
	}

	query = `
		INSERT INTO exercises(name, created_by, body_part_id, equipment_id)
		VALUES(	@name, @createdBy, @bodyPartId, @equipmentId)
		RETURNING ID
	`
	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{"name": payload.ExerciseName, "createdBy": payload.UserId, "bodyPartId": bodyPartId, "equipmentId": equipmentId}).Scan(&Id)
	if err != nil {
		var pgErr *pgConn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_name_key":
				return "", myerrs.AlreadyExitsErrMaker(payload.ExerciseName)
			}
		}
		err := fmt.Errorf("error inserting exercise %v : %w\n", payload.ExerciseName, err)
		return "", myerrs.InternalServerErrMaker(err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		err := fmt.Errorf("error commiting the transaction : %w\n", err)
		return "", myerrs.InternalServerErrMaker(err)
	}

	return Id.String(), nil
}

func (d *exerciseCD) DeleteExecise(ctx context.Context, payload mappings.DeleteExercise) error {
	// if the exercise's created by is NULL -> move it to user nullified
	// else delete the exercise

	trnx, err := d.pg.Begin(ctx)
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error creating a transaction"))
	}

	defer trnx.Rollback(ctx)

	query := `SELECT 
				id, 
				created_by 
			FROM
				exercises
			WHERE 
				name = @name AND (created_by IS NULL OR created_by = @userId) `

	var exerciseId uuid.UUID
	var createdBy *int
	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{"name": payload.ExerciseName, "userId": payload.UserId}).Scan(&exerciseId, &createdBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return myerrs.BadReqErrMaker(fmt.Errorf("exericse %v doesn't exist", payload.ExerciseName))
		}
		return myerrs.InternalServerErrMaker(fmt.Errorf("error getting id and createdBy for exericse %v : %w", payload.ExerciseName, err))
	}

	if createdBy == nil {
		query := `
			INSERT INTO user_nullified(user_id, exercise_id)
			VALUES(@userId, @exerciseId)	
		`
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"userId": payload.UserId, "exerciseId": exerciseId})
		if err != nil {
			return myerrs.InternalServerErrMaker(fmt.Errorf("error inserting data into user_nullified : %w", err))
		}

		err = trnx.Commit(ctx)
		if err != nil {
			return myerrs.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
		}

		return nil
	}

	query = `
		DELETE FROM exercises
		WHERE id = @id
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"id": exerciseId})
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error deleting rows from table : %w", err))
	}
	err = trnx.Commit(ctx)
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
	}

	return nil
}
