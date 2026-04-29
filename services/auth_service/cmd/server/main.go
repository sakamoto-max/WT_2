package main

import (
	"auth_service/internal/bootstrap"
	"os"
	"auth_service/internal/env"
)

func main() {
	
	env.Load("../../.env")

	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()
}

