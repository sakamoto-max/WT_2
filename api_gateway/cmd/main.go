package main

import (
	"os"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/bootstrap"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/env"
)

func main() {
	env.Load("../.env")

	app := bootstrap.NewApp(os.Getenv("HTTP_SERVER_ADDR"))
	app.Run()
}