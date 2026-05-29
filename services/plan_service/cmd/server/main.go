package main

import (
	"os"
	"plan_service/internal/bootstrap"
	"plan_service/internal/env"
)

func main() {

	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../../.env")
	}

	env.LookupForApi()

	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()

}
