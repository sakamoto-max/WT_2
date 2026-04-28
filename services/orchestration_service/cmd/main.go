package main

import (
	"context"
	"orchestration_service/internal/consumer"
	"orchestration_service/internal/env"
	"orchestration_service/internal/producer"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/types"
	"orchestration_service/internal/workers"
	"os"
	"os/signal"
	"sync"
	"time"
	mq "github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2-pkg/enum"
	"github.com/sakamoto-max/wt_2-pkg/logger"
	"go.uber.org/zap"
)

const NumberOfWorkers = 5

func main() {


	env.Load("../.env")

	logger := logger.NewLogger()
	defer logger.Log.Sync()

	Db, err := repository.NewDBs(logger)
	if err != nil {
		logger.Log.Fatal(err)
	}

	// rabbit mq

	conn := mq.NewConn()

	planQueue := mq.NewMessageQueue(conn, string(enum.PlanQueue))

	emailQueue := mq.NewMessageQueue(conn, string(enum.EmailQueue))

	resultQueue := mq.NewMessageQueue(conn, string(enum.ResultQueue))

	// job & workers

	jobs := make(chan types.Data, NumberOfWorkers*2)

	workers := workers.MakeWorkers(NumberOfWorkers, planQueue, emailQueue, Db, jobs, logger)

	ctx, cancel := context.WithCancel(context.Background())

	var workersWg sync.WaitGroup

	for _, worker := range workers {
		workersWg.Add(1)
		go worker.Work(ctx, &workersWg)
	}

	
	// consumer
	
	consumer := consumer.NewConsumer(Db, resultQueue, logger)
	
	msgs := consumer.GetData()
	
	go consumer.Operate(ctx, msgs)
	
	// producer
	
	targetServices := []string{string(enum.AuthService), string(enum.TrackerService)}
	
	producer := producer.NewProducer(Db, planQueue, emailQueue, logger)

	ticker := time.NewTicker(time.Second * 5)
	
	go producer.Start(ctx, &workersWg, ticker.C, jobs, &targetServices)

	// shutdown

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	logger.Log.Infow("shutdown signal received", zap.String("signal", sig.String()))

	ticker.Stop()
	cancel()
	close(jobs)
}


// what if a operation fails -> 
// retry 3 times -> send back to db -> fetch after some time -> try again -> 
// if retry > 5 -> move it to failed 