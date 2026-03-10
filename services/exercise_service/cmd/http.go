package main

import (

	"exercise_service/internal/handlers"
	"exercise_service/internal/repository"
	"exercise_service/internal/routes"
	"exercise_service/internal/services"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type httpServer struct {
	addr string
}

func NewhttpServer(addr string) *httpServer {
	return &httpServer{
		addr: addr,
	}
}

func (h *httpServer) Run(pool *pgxpool.Pool, client *redis.Client) {

	repo := repository.NewRepo(pool, client)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)
	r := routes.Routes(handler)

	log.Printf("exercise http service has started at %v\n", h.addr)

	if err := http.ListenAndServe(h.addr, r); err != nil{
		log.Fatalf("failed to listen to exercise http server %v", err)
	}
}
