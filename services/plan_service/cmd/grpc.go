package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	grpcclient "plan_service/grpc_client"
	"plan_service/internal/controllers"
	"plan_service/internal/database"
	"plan_service/internal/repository"
	"plan_service/internal/services"
	pb "workout-tracker/proto/shared/plan"

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

	exerClient := grpcclient.New()
	
	repo := repository.NewDBs(pool, redisClient)
	service := services.NewService(repo, exerClient.Client)
	controller := controllers.NewPlanController(service)
	
	grpcServer := grpc.NewServer()
	pb.RegisterPlanServiceServer(grpcServer, controller)

	go func() {
		log.Printf("grpc server has started at %v", os.Getenv("GRPC_SERVER_ADDR"))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("error listening to the grpc server : %v", err)
		}
	}()
	
	grpcServer.GracefulStop()

	if err := repo.Close(); err != nil{
		log.Println(err)
	}

	exerClient.Close()


	log.Println("gracefully shutdown")
}
