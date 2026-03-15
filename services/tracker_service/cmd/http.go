package main

import (
	"context"
	"log"
	"net/http"
	"time"
	grpcclient "tracker_service/grpc_client"
	"tracker_service/internal/database"
	"tracker_service/internal/handlers"
	"tracker_service/internal/repository"
	"tracker_service/internal/routes"
	"tracker_service/internal/services"

	// pb "workout-tracker/proto/shared/plan"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, client, err := database.InitializeDBs(ctx)
	if err != nil {
		log.Fatalf("error occured while initializing DBs : %v\n", err)
	}

	planClient := grpcclient.NewPlanClient().Connect()
	exerClient := grpcclient.NewExerClient().Connect()

	repo := repository.NewDBs(pool, client)
	service := services.NewService(repo, planClient, exerClient)
	handler := handlers.NewHandler(service)
	r := routes.Router(handler)

	log.Printf("tracker http service has started at %v\n", h.addr)

	if err := http.ListenAndServe(h.addr, r); err != nil {
		log.Fatalf("failed to listen to tracker http server %v", err)
	}
}
