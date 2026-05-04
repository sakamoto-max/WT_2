package mq_consumer

import (
	"os"
	"os/signal"
	"plan_service/internal/client"
	"plan_service/internal/env"

	"plan_service/internal/mq_consumer/consumer"
	"plan_service/internal/mq_consumer/types"
	"plan_service/internal/mq_consumer/worker"
	"plan_service/internal/repository"
	"sync"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

const numberOfWorkers = 5

func main() {

	env.Lookup()

	logger := logger.NewLogger()
	defer logger.Log.Sync()

	conn := mq.NewConn()
	Planqueue := mq.NewMessageQueue(conn, enum.QueueName_PLAN_QUEUE.String())
	resQueue := mq.NewMessageQueue(conn, enum.QueueName_RESULT_QUEUE.String())

	repo, err := repository.NewRepo()
	if err != nil {
		logger.Log.Fatalf("error opening the repos : %v", err)
	}

	jobs := make(chan types.Data, numberOfWorkers*2)

	exerClient := client.New()

	workers := worker.MakeWorkers(numberOfWorkers, repo, logger, jobs, resQueue, exerClient.Client)

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

	logger.Log.Infoln("consumer closed")

}

func InitConsumer(wg *sync.WaitGroup) {
	defer wg.Done()

	env.Lookup()

	logger := logger.NewLogger()
	defer logger.Log.Sync()

	conn := mq.NewConn()
	Planqueue := mq.NewMessageQueue(conn, enum.QueueName_PLAN_QUEUE.String())
	resQueue := mq.NewMessageQueue(conn, enum.QueueName_RESULT_QUEUE.String())

	repo, err := repository.NewRepo()
	if err != nil {
		logger.Log.Fatalf("error opening the repos for consumer : %v", err)
	}

	jobs := make(chan types.Data, numberOfWorkers*2)

	exerClient := client.New()

	workers := worker.MakeWorkers(numberOfWorkers, repo, logger, jobs, resQueue, exerClient.Client)

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

	logger.Log.Infoln("consumer closed")
}
