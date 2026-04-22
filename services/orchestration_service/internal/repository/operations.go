package repository

import (
	"context"
	"errors"
	"fmt"
	"orchestration_service/internal/types"
	"time"
	"wt/pkg/enum"

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
			5
	`
)

func (d *DB) FetchData(targerService string) (*[]types.Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryExecutionTime)
	defer cancel()

	trnx, err := d.CreateTrnx(ctx, targerService)
	if err != nil {
		return nil, err
	}

	defer trnx.Rollback(ctx)

	rows, err := trnx.Query(ctx, outboxQuery, pgx.NamedArgs{"status": enum.TaskNotCompleted})
	if err != nil {
		return nil, err
	}

	var Data []types.Data

	var id string
	var TargetService string
	var task string
	var status string
	var payload any
	var createdAt time.Time
	var numberOfTries *int

	var allIds []string

	for rows.Next() {
		err := rows.Scan(&id, &TargetService, &task, &status, &payload, &createdAt, &numberOfTries)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows : %w", err)
		}

		allIds = append(allIds, id)

		data := types.Data{
			Id:            id,
			TargetService: TargetService,
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
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"status": enum.TaskPending, "id" : id})
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
	case string(enum.AuthService):
		trnx, err := d.AuthPg.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("error in creating a auth transaction : %w", err)
		}

		return trnx, nil
	case string(enum.TrackerService):

		trnx, err := d.TrackerPg.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("error in creating a tracker transaction : %w", err)
		}

		return trnx, nil
	}

	return nil, nil
}
// only for auth
func (d *DB) UpdateNumberOfTries(ctx context.Context, id int) error {

	trnx, err := d.AuthPg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error creating a transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

	var numberOfTries int

	err = trnx.QueryRow(ctx, `
		SELECT number_of_tries FROM outbox
		WHERE ID = @id	
	`, pgx.NamedArgs{"id": id}).Scan(&numberOfTries)
	if err != nil {
		return fmt.Errorf("error occured while getting the number of tries : %w", err)
	}

	result, err := trnx.Exec(ctx, `
		UPDATE outbox
		SET number_of_tries = @number
		WHERE id = @id	
	`, pgx.NamedArgs{"number": numberOfTries + 1, "id": id})
	if err != nil {
		return fmt.Errorf("error occured while updatint the number of tries : %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNoRowsEffected
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error in commiting the transaction : %w", err)
	}

	return nil
}
func (d *DB) TaskCompletedUpdate(ctx context.Context, targetService string, id string) error {
	fmt.Println("target service", targetService)
	trnx, err := d.CreateTrnx(ctx, targetService)
	if err != nil {
		return err
	}

	fmt.Println("trnx created")
	
	defer trnx.Rollback(ctx)
	
	query := `
	UPDATE outbox
	SET status = @status
	WHERE id = @id	
	`
	
	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
		"status": enum.TaskCompleted,
		"id":     id,
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
func (d *DB) TaskPendingUpdate(ctx context.Context, targetService string, id string) error {
	trnx, err := d.CreateTrnx(ctx, targetService)
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
		"status": enum.TaskPending,
		"id":     id,
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
func (d *DB) TaskNotCompleted(ctx context.Context, targetService string, id string) error {
	query := `
		UPDATE outbox
		SET status = @status
		WHERE id = @id
	`

	trnx, err := d.CreateTrnx(ctx, targetService)
	if err != nil {
		return err
	}

	defer trnx.Rollback(ctx)

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"status": enum.TaskNotCompleted, "id": id})
	if err != nil {
		return fmt.Errorf("error in updating the task to not completed for id %v : %v", id, err)
	}


	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting : %w", err)
	}

	return nil
}
