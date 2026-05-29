package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type cancelRepo struct {
	pg *pgxpool.Pool
}

func (d *cancelRepo) DeleteTrackerIdInPG(ctx context.Context, trackerId string) error {

	query := `
		DELETE FROM tracker
		WHERE ID = @id	
	`
	_, err := d.pg.Exec(ctx, query, pgx.NamedArgs{
		"id": trackerId,
	})

	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting the tracker Id : %w", err))
	}

	return nil

}
