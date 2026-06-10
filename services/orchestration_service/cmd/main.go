package main

import (
	"orchestration_service/internal/config"
	"orchestration_service/internal/consumer"
	"orchestration_service/internal/fetcher"
	"orchestration_service/internal/server"
	"orchestration_service/internal/workers"
	"os"
	"os/signal"
)

func main() {

	config := config.LoadConfig()

	server := server.NewServer(config)

	go consumer.StartConsumer(server)
	go fetcher.StartFetcher(server)
	go workers.StartWorkers(server)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	server.Shutdown(sig.String())
}
