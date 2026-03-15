package main

import (
	"os"
	env "wt/pkg/shared/env"
)

func main() {

	env.Load()

	gRPCServer := NewgrpcServer(os.Getenv("GRPC_SERVER_ADDR"))
	go gRPCServer.Run()

	httpSer := NewhttpServer(os.Getenv("HTTP_SERVER_ADDR"))
	httpSer.Run()
}
