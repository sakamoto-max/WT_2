package main

import (
	// "context"
	"os"
	"os/signal"
	"plan_service/internal/mq_consumer/consumer"
	"plan_service/internal/mq_consumer/worker"
	"plan_service/internal/repository"
	"sync"
	"wt/pkg/enum"
	"wt/pkg/env"
	"wt/pkg/logger"
	mq "wt/pkg/queue"
	"wt/pkg/types"

	"go.uber.org/zap"
)

const numberOfWorkers = 5

func main() {

	env.Load("../../.env")

	logger := logger.NewLogger()
	defer logger.Log.Sync()

	conn := mq.NewConn()
	Planqueue := mq.NewMessageQueue(conn, string(enum.PlanQueue))
	resQueue := mq.NewMessageQueue(conn, string(enum.ResultQueue))

	repo, err := repository.NewRepo()
	if err != nil {
		logger.Log.Fatalf("error opening the repos : %v", err)
	}

	jobs := make(chan types.Data, numberOfWorkers*2)

	workers := worker.MakeWorkers(numberOfWorkers, repo, logger, jobs, resQueue)

	var workerWg sync.WaitGroup

	for _, worker := range workers {
		workerWg.Add(1)
		go worker.Work(&workerWg)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	logger.Log.Infoln("consumer has started")

	consumer := consumer.NewConsumer(Planqueue, logger, jobs)
	msgs := consumer.GetData()
	go consumer.PushDataToJobs(msgs)

	sig := <-sigChan
	logger.Log.Infow("signal received", zap.String("signal", sig.String()))

	conn.Close()
	workerWg.Wait()
	if err := repo.Close(); err != nil {
		logger.Log.Errorf("error closing the databases : %v", err)
	}

}

// msgs, err := Planqueue.Ch.Consume(string(enum.PlanQueue), "", true, false, false, false, nil)
// if err != nil {
// 	fmt.Printf("error occured while getting data from mq : %v", err)
// }

// log.Println("plan consumer has started")

// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// defer cancel()

// for msg := range msgs {
// 	switch msg.CorrelationId {
// 	case string(enum.EmptyPlanCrrId):
// 		data := utils.ConvertIntoJosn(&msg.Body)

// 		_, err := service.CreateEmptyPlan(context.TODO(), data.Payload["user_id"])
// 		if err != nil {

// 			log.Printf("error creating empty plan for  %v : %v", data.Payload["user_id"], err)

// 			data := &mq.TaskFailed{
// 				Id:            data.Id,
// 				TargerService: string(enum.PlanService),
// 				OriginatedBy:  string(enum.AuthService),
// 				TaskName:      data.Task,
// 				DbUpdateValue: string(enum.TaskNotCompleted),
// 			}

// 			dataInBytes, _ := utils.ConvertIntoBytes(data)

// 			err := ResQueue.Publish(ctx, dataInBytes, string(enum.ApplicationJsonType))

// 			if err != nil {
// 				log.Printf("task failed to send back : db_id : %v, originated_by : %v, err:%v", data.Id, data.OriginatedBy, err)
// 			}

// 			log.Println("data sent back to orc")
// 		}

// 		// d := queue.TaskStatus{
// 		// 	Id:            data.Id,
// 		// 	TargetService: string(enum.PlanService),
// 		// 	OriginatedBy:  string(enum.AuthService),
// 		// 	TaskName:      data.Task,
// 		// 	DbUpdateValue: string(enum.TaskCompleted),
// 		// }

// 		// dataInBytes, _ := utils.ConvertIntoBytes(d)

// 		// err = ResQueue.Publish(ctx, dataInBytes, string(enum.ApplicationJsonType))
// 		// if err != nil {
// 			// fmt.Println("err occured while sending back : %w", err)
// 			// outbox this
// 		// }

// 		log.Println("data sent back to orc")
// 	}

// }

// <-forever
