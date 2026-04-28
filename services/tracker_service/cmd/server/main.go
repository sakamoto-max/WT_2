package main

import (
	"os"
	"tracker_service/internal/bootstrap"
	"tracker_service/internal/env"
	// env "wt/pkg/env"
	// "github.com/sakamoto-max/wt_2-pkg/env"
)

func main() {
	env.Load("../../.env")
	
	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()
}
