package main

import (
	"plan_service/internal/bootstrap"
	"plan_service/internal/config"
)

func main() {

	// stage := os.Getenv("STAGE")
	// if stage == "" {
	// 	env.Load("../../.env")
	// }

	// env.LookupForApi()

	config := config.LoadConfig()

	app := bootstrap.NewApp(config)
	app.Run()

}
