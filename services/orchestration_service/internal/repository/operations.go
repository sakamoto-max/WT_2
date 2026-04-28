package repository

import (
	"context"
	"errors"
	"fmt"
	"orchestration_service/internal/types"
	"time"
	"github.com/jackc/pgx/v5"
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

func (d *DB) FetchData(ctx context.Context, targerService string) (*[]types.Data, error) {

	trnx, err := d.CreateTrnx(ctx, targerService)
	if err != nil {
		return nil, err
	}

	defer trnx.Rollback(ctx)

	rows, err := trnx.Query(ctx, outboxQuery, pgx.NamedArgs{"status": types.TaskNotCompleted})
	if err != nil {
		return nil, err
	}

	var Data []types.Data

	var dbId string
	var TargetService string
	var CreatedBy string
	var task string
	var status string
	var payload any
	var createdAt time.Time
	var numberOfTries int

	var allIds []string

	for rows.Next() {
		err := rows.Scan(&dbId, &TargetService, &CreatedBy, &task, &status, &payload, &createdAt, &numberOfTries)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows : %w", err)
		}

		allIds = append(allIds, dbId)

		data := types.Data{
			DbId:          dbId,
			TargetService: TargetService,
			CreatedBy: CreatedBy,
			Task:          task,
			Status:        status,
			Payload:       payload,
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
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"status": types.TaskPending, "id": id})
		if err != nil {
			return nil, fmt.Errorf("error in updating the task status to pending : %w", err)
		}
	}

	if len(Data) == 0 {
		return nil, ErrNoRowsFound
	}

	// fmt.Println(Data)

	rows.Close()

	err = trnx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in commiting the transaction : %w", err)
	}

	return &Data, nil
}
func (d *DB) CreateTrnx(ctx context.Context, targerService string) (pgx.Tx, error) {

	switch targerService {
	case string(types.AuthService):
		trnx, err := d.AuthPg.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("error in creating a auth transaction : %w", err)
		}

		return trnx, nil
	case string(types.TrackerService):

		trnx, err := d.TrackerPg.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("error in creating a tracker transaction : %w", err)
		}

		return trnx, nil
	}

	return nil, nil
}
func (d *DB) TaskCompletedUpdate(ctx context.Context, targetDbName string, dbIndex string) error {
	trnx, err := d.CreateTrnx(ctx, targetDbName)
	if err != nil {
		return err
	}


	defer trnx.Rollback(ctx)

	query := `
	UPDATE outbox
	SET status = @status
	WHERE id = @id	
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
		"status": types.TaskCompleted,
		"id":     dbIndex,
	})

	if err != nil {
		return fmt.Errorf("error updating the task to completed : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting the transaction : %w", err)
	}
	fmt.Println("commit completed")

	return nil
}
func (d *DB) TaskPendingUpdate(ctx context.Context, targetDbName string, dbIndex string) error {
	trnx, err := d.CreateTrnx(ctx, targetDbName)
	if err != nil {
		return err
	}

	defer trnx.Rollback(ctx)

	query := `
		UPDATE outbox
		SET status = @status
		WHERE id = @id	
	`
	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
		"status": types.TaskPending,
		"id":     dbIndex,
	})

	if err != nil {
		return fmt.Errorf("error updating the task status to pending : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting : %w", err)
	}

	return nil
}
func (d *DB) TaskNotCompletedUpdateTries(ctx context.Context, targetDbName string, dbIndex string) error {
	query := `
		UPDATE outbox
		SET status = @status
		WHERE id = @id
	`

	trnx, err := d.CreateTrnx(ctx, targetDbName)
	if err != nil {
		return err
	}

	defer trnx.Rollback(ctx)

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"status": types.TaskNotCompleted, "id": dbIndex})
	if err != nil {
		return fmt.Errorf("error in updating the task to not completed for id %v : %v", dbIndex, err)
	}

	var numberOfTries int

	query = `
		SELECT number_of_tries FROM outbox
		WHERE ID = @id	
	`

	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{"id": dbIndex}).Scan(&numberOfTries)
	if err != nil {
		return fmt.Errorf("error occured while getting the number of tries : %w", err)
	}

	query = `
		UPDATE outbox
		SET number_of_tries = @number
		WHERE id = @id	
	`

	_, err = trnx.Exec(ctx, query , pgx.NamedArgs{"number": numberOfTries + 1, "id": dbIndex})
	if err != nil {
		return fmt.Errorf("error occured while updatint the number of tries : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting : %w", err)
	}

	return nil
}

func (d *DB) TaskFailed(ctx context.Context, targetDbName string, dbIndex string) error {
	trnx, err := d.CreateTrnx(ctx, targetDbName)
	if err != nil {
		return err
	}

	defer trnx.Rollback(ctx)

	query := `
		UPDATE OUTBOX
		SET status = @status
		WHERE id = @id	
	`
	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"status": types.TaskFailed, "id": dbIndex})
	if err != nil {
		return fmt.Errorf("failed to update the task to failed : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to update the task to failed : %w", err)
	}

	return nil
}
