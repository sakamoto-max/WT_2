package main

import (
	"auth_service/internal/database"
	grpcclient "auth_service/internal/grpc_client"
	"auth_service/internal/handlers"
	"auth_service/internal/repository"
	"auth_service/internal/routes"
	"auth_service/internal/services"
	"context"
	"log"
	"net/http"
	"time"
)

type httpServer struct {
	addr string
}

func NewhttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func (h *httpServer) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, redisClient, err := database.InitializeDBs(ctx)
	if err != nil {
		log.Fatalf("error opening the dbs for auth http server : %v", err)
	}

	planClient := grpcclient.NewPlanClient().Connect()

	repo := repository.NewRepo(pool, redisClient)
	service := services.NewService(repo, planClient)
	handler := handlers.NewHandler(service)
	r := routes.Router(handler)

	log.Printf("auth http service has started at %v\n", h.addr)

	if err := http.ListenAndServe(h.addr, r); err != nil {
		log.Fatalf("failed to listen to auth http server %v", err)
	}
}
