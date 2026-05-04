package main

import (
	"os"
	"plan_service/internal/bootstrap"
	"plan_service/internal/env"
	// mqconsumer "plan_service/internal/mq_consumer"
	// "plan_service/internal/mq_consumer/consumer"
)

func main() {
	
	env.Lookup()

	app := bootstrap.NewApp(os.Getenv("GRPC_SERVER_ADDR"))
	app.Run()

}
