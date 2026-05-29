package main

import (
	"auth_service/internal/bootstrap"
	"auth_service/internal/env"
	"os"
)

func main() {

	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../../.env")
	}
	env.Validate()

	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()
}
