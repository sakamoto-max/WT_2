package main

import (
	"context"
	"errors"
	"orchestration_service/consumer"
	"orchestration_service/producer"
	"orchestration_service/repository"
	"orchestration_service/types"
	"orchestration_service/workers"
	"os"
	"os/signal"
	"sync"
	"time"
	"wt/pkg/enum"
	"wt/pkg/env"
	"wt/pkg/logger"
	mq "wt/pkg/queue"
	"go.uber.org/zap"

	amqp "github.com/rabbitmq/amqp091-go"
)

const NumberOfWorkers = 5

func main() {
	env.LoadNoLookUp(".env")

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
		go worker.DoWork(ctx, &workersWg)
	}
	ticker := time.NewTicker(time.Second * 10)

	// consumer

	consumer := consumer.NewConsumer(Db, resultQueue, logger)

	msgs := consumer.GetData()

	go consumer.Operate(ctx, msgs)

	// producer

	producer := producer.NewProducer(Db, planQueue, emailQueue, logger)
	go producer.Start(ctx, &workersWg, ticker.C, jobs)

	// shutdown

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	logger.Log.Infow("shutdown signal received : %v", zap.String("signal", sig.String()))

	ticker.Stop()
	cancel()
	close(jobs)
}

func ExponentialBackoff(targetFunc func(*[]byte, *mq.MessageQueue) error, a *[]byte, b *mq.MessageQueue) error {

	time.Sleep(time.Millisecond * 100)

	err := targetFunc(a, b)
	if err != nil {
		if err == amqp.ErrClosed {
			time.Sleep(time.Millisecond * 200)
			err := targetFunc(a, b)
			if err != nil {
				if errors.Is(err, amqp.ErrClosed) {
					time.Sleep(time.Millisecond * 300)
					err := targetFunc(a, b)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
