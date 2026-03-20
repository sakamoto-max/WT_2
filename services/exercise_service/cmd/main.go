package main

import (
	"os"
	env "wt/pkg/shared/env"
)

func main() {

	env.Load()

	// httpSer := NewhttpServer(os.Getenv("HTTP_SERVER_ADDR"))
	// go httpSer.Run()

	grpcSer := NewgrpcServer(os.Getenv("GRPC_SERVER_ADDR"))
	grpcSer.Run()

}
