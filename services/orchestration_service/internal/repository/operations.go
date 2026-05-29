package repository

import (
	"context"
	"errors"
	"fmt"
	"orchestration_service/internal/types"
	"orchestration_service/internal/utils"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
)

var (
	ErrNoRowsFound    = errors.New("no rows found")
	ErrNoRowsEffected = errors.New("no rows effected")
)

var (
	outboxQuery string = `
		SELECT 
			id, 
			target_service, 
			created_by,
			task, 
			status, 
			payload, 
			created_at, 
			number_of_tries 
		FROM 
			outbox
		WHERE 
			status = @status
		LIMIT 
			100
	`
)

func (a *database) FetchData(ctx context.Context, wg *sync.WaitGroup, dataChan chan<- *[]types.Data) {

	defer wg.Done()

	trnx, err := a.pg.Begin(ctx)
	if err != nil {
		dataChan <- &[]types.Data{
			{
				Err:         fmt.Errorf("failed to create a transaction : %w", err),
				ServiceName: a.dbName,
			},
		}
		return
	}

	defer trnx.Rollback(ctx)

	rows, err := trnx.Query(ctx, outboxQuery, pgx.NamedArgs{
		"status": enum.TaskStatus_TASK_NOT_COMPLETED,
	})
	if err != nil {
		dataChan <- &[]types.Data{
			{
				Err:         fmt.Errorf("failed to fetch rows from auth pg : %w", err),
				ServiceName: a.dbName,
			},
		}

		return
	}

	var Data []types.Data

	var dbId string
	var TargetService string
	var CreatedBy string
	var task string
	var status string
	var payload []byte
	var createdAt time.Time
	var numberOfTries int

	var allIds []string

	for rows.Next() {
		err := rows.Scan(&dbId, &TargetService, &CreatedBy, &task, &status, &payload, &createdAt, &numberOfTries)
		if err != nil {
			dataChan <- &[]types.Data{
				{
					Err:         fmt.Errorf("error scanning rows : %w", err),
					ServiceName: a.dbName,
				},
			}
			return
		}

		allIds = append(allIds, dbId)

		DataInJson, err := utils.ConvertToJson(payload)
		if err != nil {
			dataChan <- &[]types.Data{
				{
					Err:         fmt.Errorf("failed to convert payload into json : %w", err),
					ServiceName: a.dbName,
				},
			}
			return
		}

		data := types.Data{
			DbId:          dbId,
			TargetService: TargetService,
			CreatedBy:     CreatedBy,
			Task:          task,
			Status:        status,
			Payload:       DataInJson,
			CreatedAt:     createdAt,
			NumberOfTries: numberOfTries,
		}

		Data = append(Data, data)
	}

	query := `
		UPDATE outbox
		SET
			status = @status
		WHERE 
			id = @id
	`
	for _, id := range allIds {
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"status": enum.TaskStatus_TASK_PENDING.String(), "id": id})
		if err != nil {
			dataChan <- &[]types.Data{
				{
					Err:         fmt.Errorf("error in updating the task status to pending : %w", err),
					ServiceName: a.dbName,
				},
			}
			return
		}
	}

	if len(Data) == 0 {
		dataChan <- &[]types.Data{
			{
				ServiceName: a.dbName,
				NoData:      true,
			},
		}
		return
	}

	rows.Close()

	err = trnx.Commit(ctx)
	if err != nil {

		dataChan <- &[]types.Data{
			{
				Err:         fmt.Errorf("error in commiting the transaction : %w", err),
				ServiceName: a.dbName,
			},
		}
		return
	}

	dataChan <- &Data
}
func (a *database) UpdateTaskStatus(ctx context.Context, dbIndex string, updateValue string) error {
	query := `	
		UPDATE 
			outbox
		SET 
			status = @status
		WHERE 
			id = @id
	`
	_, err := a.pg.Exec(ctx, query, pgx.NamedArgs{
		"status": updateValue,
		"id":     dbIndex,
	})
	if err != nil {
		return fmt.Errorf("failed to update outbox task status : %w", err)
	}
	return nil
}
func (a *database) UpdateTaskStatusWithNumberOfTries(ctx context.Context, dbIndex string, updateValue string) error {

	query := `
		UPDATE 
			outbox
		SET 
			status = @STATUS,
			number_of_tries = (
				SELECT 
					number_of_tries 
				FROM 
					outbox
				WHERE 
					ID = @ID
			)+1
		WHERE 
			ID = @ID 
	`

	_, err := a.pg.Exec(ctx, query, pgx.NamedArgs{
		"STATUS": updateValue,
		"ID":     dbIndex,
	})
	if err != nil {
		return fmt.Errorf("failed to update the task status and number of tries : %w", err)
	}

	return nil
}
