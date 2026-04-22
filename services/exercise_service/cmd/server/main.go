package main

import (
	"exercise_service/internal/bootstrap"
	"os"
	env "wt/pkg/env"
)

func main() {

	env.Load("../../.env")
	
	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()
}
