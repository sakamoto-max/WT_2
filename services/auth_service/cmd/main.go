package main

import (
	"auth_service/internal/database"
	"auth_service/internal/handlers"
	"auth_service/internal/repository"
	"auth_service/internal/routes"
	"auth_service/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error occured while loading env file : %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	pool, client := database.OpenDBs(ctx)
	defer database.CloseDBs(pool, client)

	err = database.RunMigrations()
	if err != nil{
		log.Fatalf("error running migrations :%v\n", err)
	}

	repo := repository.NewRepo(pool, client)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)
	r := routes.Router(handler)
	fmt.Printf("auth server has started at %v\n", os.Getenv("SERVER_ADDR"))
	http.ListenAndServe(os.Getenv("SERVER_ADDR"), r)
	// rabbit mq
	// mq := broker.NewRabbitMq()

	// conn, err := mq.OpenMqConn()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer conn.Close()

	// channel, err := conn.Channel()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// dependency injection

}
