package repository

import (
	"context"
	"errors"
	"fmt"
	"orchestration_service/types"
	"time"
	"wt/pkg/enum"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNoRowsFound    = errors.New("no rows found")
	ErrNoRowsEffected = errors.New("no rows effected")
)

func (d *DB) FetchDataFromAuth() (*[]types.Data, error) {

	rows, err := d.AuthPg.Query(context.TODO(), `
	SELECT id, target_service, task, status, payload, created_at, number_of_tries FROM outbox
	WHERE status = @status
	LIMIT 5
	`, pgx.NamedArgs{"status": enum.TaskNotCompleted})
	if err != nil {
		return nil, fmt.Errorf("error getting data from the outbox table : %w", err)
	}
	
	defer rows.Close()

	var Data []types.Data
	
	var id string
	var TargetService string
	var task string
	var status string
	var payload any
	var createdAt time.Time
	var numberOfTries int
	
	for rows.Next() {
		err := rows.Scan(&id, &TargetService, &task, &status, &payload, &createdAt, &numberOfTries)
		if err != nil{
			return nil, fmt.Errorf("error scanning rows : %w", err)
		}

		data := types.Data{
			Id: id,
			TargetService: TargetService,
			Task: task,
			Status: status,
			Payload: payload,
			CreatedAt: createdAt,
			NumberOfTries: numberOfTries,
		}

		Data = append(Data, data)
	}

	return &Data, nil
}

func (d *DB) TaskCompletedUpdateForAuth(ctx context.Context, id string) error {
	result, err := d.AuthPg.Exec(ctx, `
		UPDATE outbox
		SET status = @status
		WHERE id = @id	
	`, pgx.NamedArgs{"status": enum.TaskCompleted, "id": id})
	if err != nil {
		return fmt.Errorf("error updating the task : %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("didnot affect any rows ")
	}

	return nil
}

func (d *DB) UpdateNumberOfTries(ctx context.Context, id int) error {

	var numberOfTries int

	trnx, err := d.AuthPg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error creating a transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

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


