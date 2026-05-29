package main

import (
	"email_service/internals/consumer"
	"email_service/internals/env"
	"email_service/internals/sender"
	"email_service/internals/server"
	"email_service/internals/worker"
	"os"
	"os/signal"
	// "github.com/sakamoto-max/rabbit_mq/types"
	// "go.uber.org/zap"
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

	server := server.NewSever()

	//	senders
	go sender.StartSenders(server)
	// workers
	go worker.StartWorkers(server)
	// consumer
	go consumer.StartConsumer(server)

	// shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan

	server.Shutdown(sig.String())
}
