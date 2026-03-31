package main

// import (
// 	// "context"
// 	// "errors"
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"orchestration_service/repository"
// 	"time"
// 	"wt/pkg/enum"
// 	"wt/pkg/env"

// 	// "log"
// 	// "orchestration_service/repository"
// 	// "os"
// 	// "os/signal"
// 	// "time"
// 	// "wt/pkg/enum"
// 	// "wt/pkg/env"
// 	mq "wt/pkg/queue"

// 	"github.com/jackc/pgx/v5"
// 	// "github.com/jackc/pgx/v5"
// 	// "github.com/redis/go-redis/v9/internal/pool"
// )

// func main() {
// 	// define a job

// 	// env.LoadNoLookUp("../.env")

// 	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	// defer cancel()

// 	// DBs, err := repository.InitializeDBs(ctx)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// conn := mq.NewConn()
// 	// defer conn.Close()

// 	// planQueue := mq.NewMessageQueue(conn, string(enum.PlanQueue))

// 	// fmt.Println("orchestration producer is on......")
// 	// // ticker := time.NewTicker(time.Second * 5)

// 	// for {
// 	// 	data, err := DBs.FetchDataFromAuth()
// 	// 	if err != nil{
// 	// 		if errors.Is(err, pgx.ErrNoRows){
// 	// 			log.Println("no rows found")
// 	// 			time.Sleep(time.Minute)
// 	// 		}

// 	// 	}
// 	// 	// fetch the data
// 	// 	// if no rows found -> print no rows found -> sleep for 1 min
// 	// 	// fetch again
// 	// 	// if rows are found -> perform the operations and push data into the queue

// 	// }

// }
