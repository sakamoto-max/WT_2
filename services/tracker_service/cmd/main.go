package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"tracker_service/internal/database"
	"tracker_service/internal/handlers"
	"tracker_service/internal/repository"
	"tracker_service/internal/routes"
	"tracker_service/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading the env files : %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, client, err := database.InitializeDBs(ctx)
	if err != nil {
		log.Fatalf("error occured while initializing DBs : %v\n", err)
	}
	defer database.CloseDBs(pool, client)

	repo := repository.NewDBs(pool, client)

	service := services.NewService(repo)

	handler := handlers.NewHandler(service)

	r := routes.Router(handler)

	http.ListenAndServe(":" + os.Getenv("SERVER_PORT"), r)
}
