package repository

import (
	"context"
	"errors"
	"fmt"

	// "plan_service/internal/domain/plan"
	"plan_service/internal/domain"
	// "plan_service/internal/mappings"

	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

var (
	ErrPlanAlreadyExists = errors.New("plan already exits")
)

type planCommandRepo struct {
	pg *pgxpool.Pool
}

func (d *planCommandRepo) CreatePlan(ctx context.Context, payload domain.CreatePlan) error {

	trnx, err := d.pg.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error creating transaction : %w\n", err))
	}

	defer trnx.Rollback(ctx)

	query := `
		INSERT INTO PLANS(USER_ID, NAME)
		VALUES (@userId, @name)	
		RETURNING ID
	`
	var planId string

	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{"userId": payload.UserId, "name": payload.PlanName}).Scan(&planId)
	if err != nil {
		var pgErr *pgConn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "one_user_one_plan_name" {
			return myerrors.BadReqErrMaker(ErrPlanAlreadyExists)
		}
		return myerrors.InternalServerErrMaker(fmt.Errorf("error inserting plan name into plans : %w\n", err))
	}

	query = `
		INSERT INTO PLAN_EXERCISES(PLAN_ID, EXERCISE_ID)
		VALUES(@planId, @exerciseId)	
	`

	for _, exerciseId := range *payload.ExerciseIds {
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"planId": planId, "exerciseId": exerciseId})
		if err != nil {
			var pgErr *pgConn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "one_plan_id_one_exercise" {
				return myerrors.BadReqErrMaker(fmt.Errorf("exercise %v is selected twice", exerciseId))
			}
			return myerrors.InternalServerErrMaker(fmt.Errorf("error inserting exercise_id %v into plan_exercises : %w", exerciseId, err))
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error commiting the trnx for creating plan : %w", err))
	}

	return nil
}


func (d *planCommandRepo) CreateEmptyPlan(ctx context.Context, userId string) error {

	query := `
		INSERT INTO plans(user_id, name)
		VALUES(@userId, @name)
	`
	_, err := d.pg.Exec(ctx, query, pgx.NamedArgs{"userId": userId, "name": "empty"})
	if err != nil {
		return fmt.Errorf("Error creating empty plan : %v", err)
	}

	return nil
}

func (d *planCommandRepo) DeletePlan(ctx context.Context, userId string, planId string) error {
	trnx, err := d.pg.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error creating a transaction : %w", err))
	}

	query := `
		DELETE FROM 
			plan_exercises
		WHERE 
			PLAN_ID = @planId
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"planId": planId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting exercises : %w", err))
	}

	query = `
		DELETE FROM 
			PLANS
		WHERE 
			ID = @id
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"id": planId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting plan : %w", err))
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
	}

	return nil

}
