package main

import (
	"os"
	"tracker_service/internal/bootstrap"
	env "wt/pkg/env"
)

func main() {

	env.Load("../../.env")
	
	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()
}
