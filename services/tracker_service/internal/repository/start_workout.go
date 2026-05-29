package repository

import (
	"context"
	"fmt"
	"tracker_service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type startWorkoutRepo struct {
	pg *pgxpool.Pool
}

func (d *startWorkoutRepo) StartWorkout(ctx context.Context, payload domain.StartWorkout) (string, error) {
	var trackerId string

	query := `
		INSERT INTO tracker(user_id, plan_id, created_at)
		VALUES(@userId, @planId, NOW())
		RETURNING id	
	`

	err := d.pg.QueryRow(ctx, query, pgx.NamedArgs{
		"userId": payload.UserId,
		"planId": payload.PlanId,
	}).Scan(&trackerId)

	if err != nil {
		return trackerId, myerrors.InternalServerErrMaker(fmt.Errorf("error starting an empty workout : %w\n", err))
	}

	return trackerId, nil
}

func (d *startWorkoutRepo) RevertStartWorkout(ctx context.Context, trackerId string) error {

	query := `
		DELETE FROM TRACKER 
		WHERE ID = @id
	`
	_, err := d.pg.Exec(ctx, query, pgx.NamedArgs{"id": trackerId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error reverting start workout : %w\n", err))
	}

	return nil
}
