package main

import (
	"context"
	"exercise_service/internal/database"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading the env file : %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	pool, client, err := database.InitializeDBs(ctx)
	if err != nil {
		log.Fatalf("error opening databases : %v\n", err)
	}

	httpSer := NewhttpServer(os.Getenv("HTTP_SERVER_ADDR"))
	go httpSer.Run(pool, client)

	grpcSer := NewgrpcServer(os.Getenv("GRPC_SERVER_ADDR"))
	grpcSer.Run(pool, client)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	// defer cancel()

	// pool, client, err := database.InitializeDBs(ctx)
	// if err != nil {
	// 	log.Fatalf("error opening databases : %v\n", err)
	// }

	// repo := repository.NewRepo(pool, client)
	// service := services.NewService(repo)
	// handler := handlers.NewHandler(service)

	// r := routes.Routes(handler)

	// fmt.Printf("exercise service has started at : %v\n", os.Getenv("SERVER_ADDR"))

	// http.ListenAndServe(":" + os.Getenv("SERVER_ADDR"), r)
}
