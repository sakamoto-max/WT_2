package main

import (
	"os"
	"tracker_service/internal/bootstrap"
	"tracker_service/internal/env"
)

func main() {
	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../../.env")
	}

	env.LookUp()

	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()
}
