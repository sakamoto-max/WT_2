package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"log"
	"time"

	"wt/pkg/enum"
	// "wt/pkg/env"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// create a db pool for the table
	// create a function to lock and fetch the rows from the table
	// print the details in terminal
	// update the data in the table and unlock\

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	pool, err := pgxpool.New(ctx, POSTGRES_CONN)
	if err != nil{
		log.Fatalf("error creating pool for the orchestration service")
	}

	for range time.Tick(time.Second * 3) {
		fetchData(pool)
	}

	pool.Close()
}

const POSTGRES_CONN string = "postgresql://postgres:root@localhost:5432/WT_AUTH?sslmode=disable"

func fetchData(pool *pgxpool.Pool)  {
	// log.Println("fetch started")
	
	var payload data
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	
	trnx, err := pool.Begin(ctx)
	if err != nil{
		// return fmt.Errorf("error creating a transaction ; %w", err)
		return
	}
	
	defer trnx.Rollback(ctx)
	
	err = trnx.QueryRow(ctx, `
	SELECT * FROM outbox
	WHERE status = $1
	LIMIT 1	FOR UPDATE SKIP LOCKED
	`,enum.TaskNotCompleted).Scan(&payload.Id, &payload.TargetService, &payload.Task, &payload.Status, &payload.Payload, &payload.CreatedAt)
	if err != nil{
		if err == pgx.ErrNoRows{
			log.Println("no tasks found")
			return 
		}
	}
	
	fmt.Println(payload)

	trnx.Exec(ctx, `
	UPDATE outbox
	SET status = $1
	WHERE id = $2	
	`,enum.TaskCompleted, payload.Id)
	
	
	err = trnx.Commit(ctx)
	if err != nil{
		// return fmt.Errorf("error commiting the transaction : %w", err)
		return
	}
	
}

type data struct{
	Id int `db:"id"`
	TargetService string `db:"target_service"`
	Task string `db:"task"`
	Status string `db:"status"`
	Payload any `db:"payload"`
	CreatedAt time.Time `db:"created_at"`
}