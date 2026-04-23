package main

import (
	"auth_service/internal/bootstrap"
	"os"

	// env "wt/pkg/env"
	env "github.com/sakamoto-max/wt_2-pkg/env"
)

func main() {
	
	env.Load("../../.env")

	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()
}