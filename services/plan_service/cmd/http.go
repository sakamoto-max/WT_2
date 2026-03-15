package main

import (
	"context"
	"log"
	"net/http"

	// grpcclient "plan_service/grpc_client"
	grpcclient "plan_service/grpc_client"
	"plan_service/internal/database"
	"plan_service/internal/handlers"
	"plan_service/internal/repository"
	"plan_service/internal/routes"
	"plan_service/internal/services"
	"time"
	// "github.com/jackc/pgx/v5/pgxpool"
	// "github.com/redis/go-redis/v9"
)

type httpServer struct {
	addr string
}

func NewhttpServer(addr string) *httpServer {
	return &httpServer{
		addr: addr,
	}
}

// func (h *httpServer) Run(grpcCli *grpcclient.ExerciseClient) {
func (h *httpServer) Run() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second *5)
	defer cancel()

	pool, redisClient, err := database.InitializeDBs(ctx)
	if err != nil{
		log.Fatalf("error opening the dbs for plan http server : %v", err)
	}

	exerClient := grpcclient.NewExerciseServiceClient().Connect()

	repo := repository.NewDBs(pool, redisClient)
	service := services.NewService(repo, exerClient)
	handler := handlers.NewHandler(service)
	r := routes.Router(handler)

	log.Printf("plan http service has started at %v\n", h.addr)

	if err := http.ListenAndServe(h.addr, r); err != nil {
		log.Fatalf("failed to listen to plan http server %v", err)
	}
}
