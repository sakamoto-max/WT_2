package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

var (
	ErrExerAlrExistsInPlan = errors.New("exercise already exists")
)

type DbErr struct {
	err        error
	exerciseId string
}

func (d *DbErr) Error() string {
	return d.err.Error()
}

func (d *DbErr) GetExerciseId() string {
	return d.exerciseId
}

type planExerciseRepo struct {
	pg *pgxpool.Pool
}

func (d *planExerciseRepo) AddExercisesToPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {

	trnx, err := d.pg.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error creating a transaction : %w", err))
	}

	defer trnx.Rollback(ctx)

	query := `
			INSERT INTO 
				plan_exercises(plan_id, exercise_id)
			VALUES 
				(@planId, @exerciseId)
		`

	for _, id := range *exerciseIDs {
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"planId": planId, "exerciseId": id})
		if err != nil {
			var pgErr *pgConn.PgError

			if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "one_plan_id_one_exercise" {
				return &DbErr{err: ErrExerAlrExistsInPlan, exerciseId: id}
			}

			return myerrors.InternalServerErrMaker(fmt.Errorf("error inserting the exericse with id %v : %w", id, err))
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
	}

	return nil
}

func (d *planExerciseRepo) RemoveExerciseFromPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {

	trnx, err := d.pg.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error creating a transaction : %w", err))
	}

	defer trnx.Rollback(ctx)

	query := `
		DELETE FROM 
			plan_exercises
		WHERE 
			PLAN_ID = @planId AND EXERCISE_ID = @exerciseId
	`

	for _, id := range *exerciseIDs {
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"planId": planId, "exerciseId": id})
		if err != nil {
			return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting exercise with id %v : %w", id, err))
		}


	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
	}

	return nil
}
