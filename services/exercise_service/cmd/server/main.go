package main

import (
	"exercise_service/internal/bootstrap"
	"exercise_service/internal/env"
	"os"
)

func main() {

	stage := os.Getenv("STAGE")
	if stage != "" {
		env.Load("../../.env")
	}

	env.LookUp()

	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()
}
