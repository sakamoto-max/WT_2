package main

import (
	"context"
	"orchestration_service/internal/consumer"
	"orchestration_service/internal/database"
	"orchestration_service/internal/env"
	"orchestration_service/internal/fetcher"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/types"
	"orchestration_service/internal/workers"
	"os"
	"os/signal"
	"sync"
	"time"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

const NumberOfWorkers = 5

func main() {

	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../.env")
	}
	env.LookUp()

	logger := logger.NewLogger()
	defer logger.Log.Sync()

	// create all the dependencies
	// and inject

	authPool, err := database.NewPool(os.Getenv("AUTH_POSTGRES_CONN"))
	if err != nil {
		logger.Log.Fatalw("failed to connect to auth pg", zap.Error(err))
	}

	trackerPool, err := database.NewPool(os.Getenv("TRACKER_POSTGRES_CONN"))
	if err != nil {
		logger.Log.Fatalw("failed to connect to tracker pg", zap.Error(err))
	}

	authDb := repository.RegisterDb(authPool, enum.ServiceName_AUTH_SERVICE.String())
	trackerDb := repository.RegisterDb(trackerPool, enum.ServiceName_TRACKER_SERVICE.String())

	Db := repository.NewDb(authDb, trackerDb)

	// queues : planQueue, emailQueue, resultQueue
	conn := mq.NewConn()

	planQueue := mq.NewMessageQueue(conn, enum.QueueName_PLAN_QUEUE.String())

	emailQueue := mq.NewMessageQueue(conn, enum.QueueName_EMAIL_QUEUE.String())

	resultQueue := mq.NewMessageQueue(conn, enum.QueueName_RESULT_QUEUE.String())

	// jobs chan
	jobs := make(chan types.Data, NumberOfWorkers*2)

	// ctx
	ctx, cancel := context.WithCancel(context.Background())

	// start consumer
	consumer := consumer.NewConsumer(resultQueue, logger, jobs)
	go consumer.StartListening()

	// start producer
	ticker := time.NewTicker(time.Second * 30)

	targetServices := []string{enum.ServiceName_AUTH_SERVICE.String(), enum.ServiceName_TRACKER_SERVICE.String()}

	var producerWg sync.WaitGroup

	fetcher := fetcher.NewFetcher(&Db, logger, &targetServices, jobs, ticker.C)

	producerWg.Add(1)

	go fetcher.Start(ctx, &producerWg)

	// start workers

	workers := workers.MakeWorkers(NumberOfWorkers, planQueue, emailQueue, &Db, jobs, logger)

	var workersWg sync.WaitGroup

	for _, worker := range workers {
		workersWg.Add(1)
		go worker.Work(ctx, &workersWg)
	}
	logger.Log.Infow("workers have started", zap.Int("number of workers", NumberOfWorkers))

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	logger.Log.Infow("shutdown signal received", zap.String("signal", sig.String()))

	ticker.Stop() // stops producer
	cancel()
	producerWg.Wait()
	logger.Log.Infoln("producer have stopped")

	conn.Close() // stops consumer
	logger.Log.Infoln("consumer has closed")

	close(jobs)
	workersWg.Wait()
	logger.Log.Infoln("workers have stopped")
}

// what if a operation fails ->
// retry 3 times -> send back to db -> fetch after some time -> try again ->
// if retry > 5 -> move it to failed
