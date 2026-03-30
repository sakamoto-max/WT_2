package main

import (
	"os"
	env "wt/pkg/env"
)

func main() {

	env.Load("../.env")

	gRPCServer := NewgrpcServer(os.Getenv("GRPC_SERVER_ADDR"))
	gRPCServer.Run()


}
