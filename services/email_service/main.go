package main

import (
	"email_service/consumer"
	"email_service/producer"
	"email_service/services"
	"email_service/worker"
	// "log"
	"os"
	"os/signal"
	"sync"
	"wt/pkg/enum"
	"wt/pkg/logger"
	"wt/pkg/queue"
	"wt/pkg/types"
	// "wt/pkg/utils"

	"go.uber.org/zap"
)

const NumberOfWorkers = 5

func main() {

	conn := queue.NewConn()

	emailQueue := queue.NewMessageQueue(conn, string(enum.EmailQueue))

	resQueue := queue.NewMessageQueue(conn, string(enum.ResultQueue))

	logger := logger.NewLogger()
	defer logger.Log.Sync()

	service := services.NewService(logger)
	jobs := make(chan types.Data, NumberOfWorkers * 2)

	producer := producer.NewProducer(logger, resQueue)

	workers := worker.MakeWorkers(NumberOfWorkers, logger, service, jobs, producer)

	var workerWg sync.WaitGroup

	for _, worker := range workers{
		workerWg.Add(1)
		go worker.Work(&workerWg)
	}




	consumer := consumer.NewConsumer(emailQueue, logger, jobs)
	msgs := consumer.Start()
	go consumer.PushToJobs(msgs)

	// msgs, err := emailQueue.Consume(string(enum.EmailQueue))
	// if err != nil {
	// 	log.Fatalf("error consuming from the email queue : %v", err)
	// }

	// log.Printf("email consumer has started")

	// for msg := range msgs {
	// 	data := utils.ConvertIntoJosn(&msg.Body)
	// 	email := data.Payload["email"]
	// 	log.Printf("sending email to : %v", email)
	// 	continue
	// }


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <- sigChan

	logger.Log.Infow("shutdown signal received", zap.String("signal", sig.String()))

	conn.Close()
	workerWg.Wait()


}