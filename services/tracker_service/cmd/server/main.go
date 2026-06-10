package main

import (
	"tracker_service/internal/bootstrap"
	"tracker_service/internal/config"
)

func main() {
	config := config.LoadConfig()

	app := bootstrap.NewApp(config)
	app.Run()
}
