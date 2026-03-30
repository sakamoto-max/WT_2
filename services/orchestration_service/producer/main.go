package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"wt/pkg/enum"
	mq "wt/pkg/queue"
	"wt/pkg/utils"
	"wt/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// create a db pool for the table
	// create a function to fetch one row at a time
	// create a messging queue
	// send the data into the queue
	// get id back -> update the row to completed

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, err := pgxpool.New(ctx, POSTGRES_CONN)
	if err != nil {
		log.Fatalf("error creating pool for the orchestration service")
	}

	defer pool.Close()

	conn := mq.NewConn()
	defer conn.Close()



	var forever chan int

	fmt.Println("orchestration producer is on......")
	for {
		time.Sleep(time.Second * 3)
		data, err := fetchData(pool)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				log.Println("no rows found")
			}
		}else{
			fmt.Println("data :", data)
			dataInBytes, _ := utils.ConvertIntoBytes(data)
		
			planQueue := mq.NewMessageQueue(conn, string(enum.PlanQueue))
		
			err = planQueue.Publish(ctx, dataInBytes, string(enum.ApplicationJsonType))
			if err != nil{
				log.Println(err)
			}
		}
	}


	// go func() {

	// 	data, err := fetchData(pool)
	// 	if err != nil {
	// 		if errors.Is(err, ErrNoRowsFound) {
	// 			log.Println("no rows found")
	// 			return
	// 		}
	// 		fmt.Println("error occured :", err)
	// 	}
	
	// 	fmt.Println("data :", data)
	// 	dataInBytes, _ := utils.ConvertIntoBytes(data)
	
	// 	planQueue := mq.NewMessageQueue(conn, string(enum.PlanQueue))
	
	// 	err = planQueue.Publish(ctx, dataInBytes, string(enum.ApplicationJsonType))
	// 	if err != nil{
	// 		log.Println(err)
	// 	}

	// }()

	<- forever
}

const POSTGRES_CONN string = "postgresql://postgres:root@localhost:5432/WT_AUTH?sslmode=disable"

func fetchData(pool *pgxpool.Pool) (*types.Data, error) {
	// log.Println("fetch started")

	var payload types.Data

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	trnx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating a transaction ; %w", err)
	}

	defer trnx.Rollback(ctx)

	err = trnx.QueryRow(ctx, `
	SELECT id, target_service, task, status, payload, created_at, number_of_tries FROM outbox
	WHERE status = @status
	LIMIT 1
	`, pgx.NamedArgs{"status": enum.TaskNotCompleted}).Scan(&payload.Id, &payload.TargetService, &payload.Task, &payload.Status, &payload.Payload, &payload.CreatedAt, &payload.NumberOfTries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoRowsFound
		}

		return nil, fmt.Errorf("error getting data from the outbox table : %w", err)
	}

	_, err = trnx.Exec(ctx, `
	UPDATE outbox
	SET status = @status
	WHERE id = @id
	`, pgx.NamedArgs{"status": enum.TaskPending, "id": payload.Id})
	if err != nil {
		return &payload, fmt.Errorf("error updating the outbox table rows : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return &payload, fmt.Errorf("error commiting the data : %w", err)
	}

	return &payload, nil
}

func FinishTask(ctx context.Context, pool *pgxpool.Pool, id int) error {
	result, err := pool.Exec(ctx, `
		UPDATE outbox
		SET status = @status
		WHERE id = @id	
	`, pgx.NamedArgs{"status": enum.TaskCompleted, "id": id})
	if err != nil {
		return fmt.Errorf("error finishing the task : %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("didnot affect any rows ")
	}

	return nil
}

func UpdateNumberOfTries(ctx context.Context, pool *pgxpool.Pool, id int) error {

	var numberOfTries int

	trnx, err := pool.Begin(ctx)
	if err != nil{
		return fmt.Errorf("error creating a transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

	err = trnx.QueryRow(ctx, `
		SELECT number_of_tries FROM outbox
		WHERE ID = @id	
	`,pgx.NamedArgs{"id" : id}).Scan(&numberOfTries)
	if err != nil{
		return fmt.Errorf("error occured while getting the number of tries : %w", err)
	}

	result, err := trnx.Exec(ctx, `
		UPDATE outbox
		SET number_of_tries = @number
		WHERE id = @id	
	`, pgx.NamedArgs{"number" : numberOfTries + 1, "id" : id})
	if err != nil{
		return fmt.Errorf("error occured while updatint the number of tries : %w", err)
	}

	if result.RowsAffected() == 0{
		return ErrNoRowsEffected 
	}

	err = trnx.Commit(ctx)
	if err != nil{
		return fmt.Errorf("error in commiting the transaction : %w", err)
	}

	return nil
}

var (
	ErrNoRowsFound = errors.New("no rows found")
	ErrNoRowsEffected = errors.New("no rows effected")
)
