package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
)

type queueDb struct {
	pg *pgxpool.Pool
}

const (
	queryExecutionTime = time.Second * 5
)

func (q *queueDb) Insert(data mqTypes.Data) error {

	ctx, cancel := context.WithTimeout(context.Background(), queryExecutionTime)
	defer cancel()

	query := `
		INSERT INTO 
			push_to_queue_failed(
				target_service,
				target_db_id,
				task_name,
				status,
				number_of_tries
			) 

		VALUES(
			@targetService,
			@targetDbId,
			@taskName,
			@status,
			@numberOfTries
		)
	`

	_, err := q.pg.Exec(ctx, query, pgx.NamedArgs{
		"targetDbId":    data.DbId,
		"targetService": data.TargetService,
		"numberOfTries": 4,
		"status":        data.TaskStatus,
		"taskName":      data.TaskName,
	})
	if err != nil {
		return fmt.Errorf("failed to insert into db : %w", err)
	}

	return nil
}

func (q *queueDb) Fetch() (*[]mqTypes.Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryExecutionTime)
	defer cancel()

	query := `
		SELECT 
			target_service,
			target_db_id,
			task_name,
			status
		FROM 
			push_to_queue_failed
		LIMIT 
			100	
	`

	rows, err := q.pg.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from the push_to_queue_failed")
	}

	var allData []mqTypes.Data

	var targetService string
	var targetDbId string
	var taskName string
	var status string

	for rows.Next() {
		err := rows.Scan(&targetService, &targetDbId, &taskName, &status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan the data from push_to_queue_failed")
		}

		data := mqTypes.Data{
			DbId:          targetDbId,
			TaskName:      taskName,
			TargetService: targetService,
			TaskStatus: status,
			SentBy: enum.ServiceName_PLAN_SERVICE.String(),
		}

		allData = append(allData, data)
	}


	if len(allData) == 0 {
		return nil, nil
	}

	return &allData, nil
}
