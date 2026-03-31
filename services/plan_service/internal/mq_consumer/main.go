package main

import (
	"context"
	"fmt"
	"log"
	"plan_service/internal/database"
	"plan_service/internal/repository"
	"plan_service/internal/services"
	"time"
	"wt/pkg/enum"
	mq "wt/pkg/queue"
	"wt/pkg/env"
	"wt/pkg/utils"
)

func main() {
	// connect to rabbit mq
	// create a channel
	// declare the queue
	// keep reading data forever
	// if data is received
	// perform the operation
	// return failure or success
	env.Load("../../.env")

	conn := mq.NewConn()
	queue := mq.NewMessageQueue(conn, string(enum.PlanQueue))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, redisClient, err := database.InitializeDBs(ctx)
	if err != nil {
		log.Fatalf("error opening the dbs for plan http server : %v", err)
	}

	repo := repository.NewDBs(pool, redisClient)
	service := services.NewPlanService(repo)

	var forever chan int

	log.Println("consumer is listening.....")
	go func() {
		
		msgs, err := queue.Ch.Consume(string(enum.PlanQueue), "", true, false, false, false, nil)
		if err != nil {
			fmt.Printf("error occured while getting data from mq : %v", err)
		}

		for msg := range msgs {
			switch msg.CorrelationId {
			case string(enum.EmptyPlanCrrId):
				data := utils.ConvertIntoJosn(&msg.Body)
				log.Printf("message recieved : task %v : user_id : %v", data.Task, data.Payload["user_id"])
				_, err := service.CreateEmptyPlan(context.TODO(), data.Payload["user_id"])
				if err != nil {
					fmt.Printf("error creating empty plan for  %v : %v", data.Id, err)
				}

			}
		}
	}()

	<-forever
}

