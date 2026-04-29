package main

import (
	"exercise_service/internal/bootstrap"
	"exercise_service/internal/env"
	"os"
)

func main() {

	env.Load("../../.env")
	
	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()
}
