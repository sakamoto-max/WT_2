package main

import (
	"auth_service/internal/bootstrap"
	"auth_service/internal/config"
)

func main() {

	config := config.LoadConfig()

	app := bootstrap.NewApp(config)
	app.Run()
}
