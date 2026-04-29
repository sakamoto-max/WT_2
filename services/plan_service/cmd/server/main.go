package main

import (
	"os"
	"plan_service/internal/bootstrap"
	"plan_service/internal/env"
)

func main() {
	
	env.Load("../../.env")

	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()

}
