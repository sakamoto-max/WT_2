package main

import (
	"plan_service/internal/bootstrap"
	"plan_service/internal/config"
)

func main() {

	config := config.LoadConfig()

	app := bootstrap.NewApp(config)
	app.Run()

}
