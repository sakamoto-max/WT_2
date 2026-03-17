package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	grpcclient "tracker_service/grpc_client"
	"tracker_service/internal/controllers"
	"tracker_service/internal/database"
	"tracker_service/internal/repository"
	"tracker_service/internal/services"
	trackerpb "workout-tracker/proto/shared/tracker"

	"google.golang.org/grpc"
)

type grpcServer struct {
	addr string
}

func NewgrpcServer(addr string) *grpcServer {
	return &grpcServer{
		addr: addr,
	}
}

func (g *grpcServer) Run() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, redisClient, err := database.InitializeDBs(ctx)
	if err != nil {
		log.Fatalf("error opening the dbs for plan grpc server: %v", err)
	}

	lis, err := net.Listen("tcp", os.Getenv("GRPC_SERVER_ADDR"))

	if err != nil {
		log.Fatalf("failed to listen to tcp : %v", err)
	}

	grpcServer := grpc.NewServer()
	exerClient := grpcclient.NewExerClient().Connect()
	planClient := grpcclient.NewPlanClient().Connect()

	repo := repository.NewDBs(pool, redisClient)
	service := services.NewService(repo, planClient, exerClient)
	controller := controllers.NewTrackerController(service)

	trackerpb.RegisterTrackerServiceServer(grpcServer, controller)
	
	log.Printf("grpc server has started at %v", os.Getenv("GRPC_SERVER_ADDR"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error listening to the grpc server : %v", err)
	}
}
