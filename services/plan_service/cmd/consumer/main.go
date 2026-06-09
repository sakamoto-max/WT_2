package main

import (
	"os"
	"os/signal"
	"plan_service/internal/mq_consumer/server"
	"plan_service/internal/env"
	"plan_service/internal/mq_consumer/consumer"
	"plan_service/internal/mq_consumer/sender"
	"plan_service/internal/mq_consumer/worker"
)

func main() {

	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../../.env")
	}

	env.LookupForConsumer()

	server := server.NewServer()

	go sender.StartSenders(server)
	go worker.StartWorkers(server)
	go consumer.StartConsumer(server)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan

	server.ShutDown(sig.String())
}
