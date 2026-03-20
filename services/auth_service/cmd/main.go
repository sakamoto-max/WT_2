package main

import (
	"os"

	env "wt/pkg/shared/env"
)

func main() {
	
	env.Load()

	grpcServer := NewgrpcServer(os.Getenv("GRPC_SERVER_ADDR"))
	grpcServer.Run()	
}

// httpServer := NewhttpServer(os.Getenv("HTTP_SERVER_ADDR"))
// go httpServer.Run()
// rabbit mq
// mq := broker.NewRabbitMq()

// conn, err := mq.OpenMqConn()
// if err != nil {
// 	log.Fatal(err)
// }

// defer conn.Close()

// channel, err := conn.Channel()
// if err != nil {
// 	log.Fatal(err)
// }

// dependency injection
