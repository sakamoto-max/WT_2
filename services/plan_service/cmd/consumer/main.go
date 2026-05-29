package main

import (
	"os"
	"os/signal"
	"plan_service/internal/config"
	"plan_service/internal/env"
	"plan_service/internal/mq_consumer/consumer"
	"plan_service/internal/mq_consumer/sender"
	"plan_service/internal/mq_consumer/worker"

	// "go.uber.org/zap"
)

// const (
// 	numberOfWorkers = 5
// 	numberOfSenders = 5
// )

func main() {

	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../../.env")
	}

	env.LookupForConsumer()

	config := config.NewConfig()

	go sender.StartSenders(config)
	go worker.StartWorkers(config)
	go consumer.StartConsumer(config)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan

	config.ShutDown(sig.String())
}
