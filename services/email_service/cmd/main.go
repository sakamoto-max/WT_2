package main

import (
	"email_service/internals/config"
	"email_service/internals/consumer"
	"email_service/internals/sender"
	"email_service/internals/server"
	"email_service/internals/worker"
	"os"
	"os/signal"
)

func main() {

	config := config.LoadConfig()
	
	server := server.NewSever(config)

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
