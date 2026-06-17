package main

import (
	"github.com/sakamoto-max/wt_2/api_gateway/internals/bootstrap"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/config"
	
)


func main() {

	config := config.LoadConfig()

	app := bootstrap.NewApp(config)
	app.Run()
}
