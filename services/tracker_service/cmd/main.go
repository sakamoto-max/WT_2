package main

import (
	"os"
	env "wt/pkg/env"
)

func main() {

	env.Load("../.env")
	
	grpcSer := NewgrpcServer(os.Getenv("GRPC_SERVER_ADDR"))
	grpcSer.Run()

}
