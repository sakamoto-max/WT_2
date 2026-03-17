package main

import (
	"os"
	env "wt/pkg/shared/env"
)

func main() {

	env.Load()
	
	grpcSer := NewgrpcServer(os.Getenv("GRPC_SERVER_ADDR"))
	go grpcSer.Run()

	httpSer := NewhttpServer(os.Getenv("HTTP_SERVER_ADDR"))
	httpSer.Run()

}
