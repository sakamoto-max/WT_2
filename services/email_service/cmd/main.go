package main

import (
	"email_service/internals/consumer"
	"email_service/internals/database"
	"email_service/internals/env"
	"email_service/internals/repostitory"
	"email_service/internals/sender"
	"email_service/internals/services"
	"email_service/internals/types"
	"email_service/internals/worker"
	"os"
	"os/signal"
	"sync"

	"github.com/sakamoto-max/rabbit_mq/queue"
	// "github.com/sakamoto-max/rabbit_mq/types"

	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

const (
	NumberOfWorkers = 5
	NumberOfSenders = 5
)

func main() {

	// env
	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../.env")
	}
	env.LookUp()

	// logger
	logger := logger.NewLogger()
	defer logger.Log.Sync()
	// db
	pool := database.NewDb(os.Getenv("POSTGRES_CONN"), logger)

	db := repostitory.RegisterDb(pool, logger)

	logger.Log.Infoln("connected to the database")
	// mq
	conn := queue.NewConn()
	logger.Log.Infoln("connected to the rabbit mq")

	emailQueue := queue.NewMessageQueue(conn, enum.QueueName_EMAIL_QUEUE.String())

	resQueue := queue.NewMessageQueue(conn, enum.QueueName_RESULT_QUEUE.String())
	// service
	service := services.NewService(logger)
	// sender chan
	senderChan := make(chan types.Data, NumberOfSenders*2)
	// jobs chan
	jobsChan := make(chan types.Data, NumberOfWorkers*2)
	// senders
	senders := sender.MakeSenders(NumberOfSenders, logger, resQueue, senderChan, db)

	var senderWg sync.WaitGroup

	for _, sender := range senders {
		senderWg.Add(1)
		go sender.Start(&senderWg)
	}
	logger.Log.Infow("senders have started", zap.Int("number of senders", NumberOfSenders))
	// workers
	workers := worker.MakeWorkers(NumberOfWorkers, logger, service, jobsChan, senderChan)

	var workerWg sync.WaitGroup

	for _, worker := range workers {
		workerWg.Add(1)
		go worker.Work(&workerWg)
	}
	logger.Log.Infow("workers have started", zap.Int("number of workers", NumberOfWorkers))
	// consumer

	consumer := consumer.NewConsumer(emailQueue, logger, jobsChan)

	go consumer.StartListening()

	logger.Log.Infoln("consumer has started")

	// shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan

	logger.Log.Infow("shutdown signal received", zap.String("signal", sig.String()))

	// close consumer -> close the emailqueue
	err := emailQueue.Ch.Close()
	if err != nil {
		logger.Log.Errorw("falied to close the email queue channel", zap.Error(err))
	}

	logger.Log.Infoln("consumer have stopped")
	// close worker
	close(jobsChan)
	workerWg.Wait()

	logger.Log.Infoln("workers have stopped")
	// close senders
	close(senderChan)
	senderWg.Wait()

	logger.Log.Infoln("senders have stopped")
	// close pool
	pool.Close()

	logger.Log.Infoln("database connection is closed")
	// mq conn
	conn.Close()
	logger.Log.Infoln("rabbit mq connection is closed")

}
