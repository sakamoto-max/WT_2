package main

import (
	"auth_service/internal/controllers"
	"auth_service/internal/database"
	grpcclient "auth_service/internal/grpc_client"
	"auth_service/internal/repository"
	"auth_service/internal/services"
	"context"
	"log"
	"net"
	"time"

	pb "workout-tracker/proto/shared/auth"

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

	lis, err := net.Listen("tcp", g.addr)

	if err != nil {
		log.Fatalf("failed to listen to tcp : %v", err)
	}

	grpcServer := grpc.NewServer()
	planClient := grpcclient.NewPlanClient().Connect()

	repo := repository.NewRepo(pool, redisClient)
	service := services.NewService(repo, planClient)
	controller := controllers.NewAuthController(service)

	pb.RegisterAuthServiceServer(grpcServer, controller)

	log.Printf("grpc server has started at %v", g.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error listening to the grpc server : %v", err)
	}
}
