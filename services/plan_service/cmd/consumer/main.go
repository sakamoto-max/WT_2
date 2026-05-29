package main

import (
	"os"
	"os/signal"
	"plan_service/internal/client"
	"plan_service/internal/database"
	"plan_service/internal/env"
	"plan_service/internal/mq_consumer/consumer"
	"plan_service/internal/mq_consumer/sender"
	"plan_service/internal/mq_consumer/types"
	"plan_service/internal/mq_consumer/worker"
	"plan_service/internal/repository"
	"sync"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

const (
	numberOfWorkers = 5
	numberOfSenders = 5
)

func main() {

	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../../.env")
	}

	env.LookupForConsumer()

	logger := logger.NewLogger()
	defer logger.Log.Sync()

	mqConn := mq.NewConn()
	Planqueue := mq.NewMessageQueue(mqConn, enum.QueueName_PLAN_QUEUE.String())
	resQueue := mq.NewMessageQueue(mqConn, enum.QueueName_RESULT_QUEUE.String())

	pool, err := database.NewPgConn()
	if err != nil {
		logger.Log.Fatalw("failed to open postgres connection for plan consumer", zap.Error(err))
	}

	repo := repository.NewDb(pool)
	logger.Log.Infoln("connected to db")

	// jobs chan

	jobs := make(chan types.Data, numberOfWorkers*2)

	// sender chan
	senderChan := make(chan mqTypes.Data, numberOfSenders*2)

	// senders

	senders := sender.MakeSenders(numberOfSenders, senderChan, resQueue, logger)

	var senderWg sync.WaitGroup

	for _, sender := range senders {
		senderWg.Add(1)
		go sender.Start(&senderWg)
	}
	logger.Log.Infow("senders have started", zap.Int("number of senders", numberOfSenders))

	exerConn := client.NewEmptyClient().OpenConnection(os.Getenv("EXERCISE_GRPC_SERVER_ADDR"))
	exerciseClient := exerConn.CreateExerciseClient()

	workers := worker.MakeWorkers(numberOfWorkers, repo, logger, jobs, senderChan, exerciseClient)

	var workerWg sync.WaitGroup

	for _, worker := range workers {
		workerWg.Add(1)
		go worker.Work(&workerWg)
	}
	logger.Log.Infow("workers have started", zap.Int("number of workers", numberOfWorkers))

	consumer := consumer.NewConsumer(Planqueue, logger, jobs)
	go consumer.Start()
	logger.Log.Infoln("consumer has started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	logger.Log.Infow("signal received", zap.String("signal", sig.String()))

	// mqConn.Close()
	Planqueue.Ch.Close()
	logger.Log.Infoln("consumer has closed")

	close(jobs)

	workerWg.Wait()
	logger.Log.Infoln("workers have stopped")

	close(senderChan)
	senderWg.Wait()
	logger.Log.Infoln("senders have stopped")

	mqConn.Close()
	
	pool.Close()
	logger.Log.Infoln("db connection  closed")

	logger.Log.Infoln("shutdown")
}
