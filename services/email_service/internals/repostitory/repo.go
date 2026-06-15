package repostitory

import (
	"context"
	"email_service/internals/types"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sakamoto-max/wt_2_pkg/logger"
)

const (
	queryExecutionTime = time.Second * 5
)

type Db struct {
	pg     *pgxpool.Pool
	logger *logger.MyLogger
}

type RepoIFace interface {
	PushToFailed(data types.Data, numberOfTries int, status string, Err error) error
}

func RegisterDb(pool *pgxpool.Pool, logger *logger.MyLogger) RepoIFace {
	return &Db{
		pg:     pool,
		logger: logger,
	}
}

func (d *Db) PushToFailed(data types.Data, numberOfTries int, status string, Err error) error {

	ctx, cancel := context.WithTimeout(context.Background(), queryExecutionTime)
	defer cancel()

	query := `
		INSERT INTO PUSH_TO_QUEUE_FAILED(target_db_id, target_service, number_of_tries, status, reason)
		VALUES(
			@dbId,
			@targetService,
			@numberOfTries,
			@status,
			@reason
		)
	`

	_, err := d.pg.Exec(ctx, query, pgx.NamedArgs{
		"dbId":          data.DbId,
		"targetService": data.TargetService,
		"numberOfTries": numberOfTries,
		"status":        status,
		"reason":        Err.Error(),
	})

	if err != nil {
		return fmt.Errorf("failed to insert data into push_to_queue_failed table : %w", err)
	}

	return nil

}
