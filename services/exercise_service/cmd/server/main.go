package main

import (
	"exercise_service/internal/bootstrap"
	"exercise_service/internal/config"
)

func main() {

	config := config.LoadConfig()

	app := bootstrap.NewApp(config)
	app.Run()
}
