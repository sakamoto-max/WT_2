package main

import (
	"context"
	"exercise_service/internal/database"
	"exercise_service/internal/handlers"
	"exercise_service/internal/repository"
	"exercise_service/internal/routes"
	"exercise_service/internal/services"
	"log"
	"net/http"
	"time"
)

type httpServer struct {
	addr string
}

func NewhttpServer(addr string) *httpServer {
	return &httpServer{
		addr: addr,
	}
}

func (h *httpServer) Run() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	pool, redisClient, err := database.InitializeDBs(ctx)
	if err != nil {
		log.Fatalf("error opening the dbs for plan http server : %v", err)
	}

	repo := repository.NewRepo(pool, redisClient)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)
	r := routes.Routes(handler)

	log.Printf("exercise http service has started at %v\n", h.addr)

	if err := http.ListenAndServe(h.addr, r); err != nil {
		log.Fatalf("failed to listen to exercise http server %v", err)
	}
}
