package main

import (
	"os"

	env "wt/pkg/env"
)

func main() {
	
	env.Load("../.env")

	grpcServer := NewgrpcServer(os.Getenv("GRPC_SERVER_ADDR"))
	grpcServer.Run()	
}