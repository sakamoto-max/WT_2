package main

import (
	"os"

	"github.com/sakamoto-max/wt_2/api_gateway/internals/bootstrap"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/env"
	// "github.com/swaggo/http-swagger" 
	// _ "github.com/sakamoto-max/wt_2/api_gateway/cmd/server/docs"
)


func main() {

	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../../.env")
	}

	env.LookUp()

	app := bootstrap.NewApp(os.Getenv("HTTP_SERVER_ADDR"))
	app.Run()
}
