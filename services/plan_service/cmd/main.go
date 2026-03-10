package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"plan_service/internal/database"
	handler "plan_service/internal/handlers"
	"plan_service/internal/repository"
	"plan_service/internal/routes"
	"plan_service/internal/services"
	"time"

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

	err = database.RunMigrationsUp()
	if err != nil{
		log.Fatalf("error occured while running migrations : %v\n", err)
	}

	repo := repository.NewDBs(pool, client)

	service := services.NewService(repo)

	handler := handler.NewHandler(service)

	r := routes.Router(handler)

	fmt.Printf("plan server is listening at %v\n", os.Getenv("SERVER_PORT"))

	err = http.ListenAndServe(":" + os.Getenv("SERVER_PORT"), r)
	if err != nil{
		log.Fatalf("error listening : %v\n", err)
	}
}
