package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"exercise_service/internal/controllers"
	"exercise_service/internal/database"
	// "exercise_service/internal/handlers"
	"exercise_service/internal/repository"
	"exercise_service/internal/services"
	exercisepb "workout-tracker/proto/shared/exercise"

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
		log.Fatalf("error opening the dbs for plan http server : %v", err)
	}

	lis, err := net.Listen("tcp", os.Getenv("GRPC_SERVER_ADDR"))

	if err != nil {
		log.Fatalf("failed to listen to tcp : %v", err)
	}

	grpcServer := grpc.NewServer()

	repo := repository.NewRepo(pool, redisClient)
	
	service := services.NewService(repo)

	controller := controllers.NewExerController(service)

	exercisepb.RegisterExerciseServiceServer(grpcServer, controller)

	log.Printf("grpc server has started at %v", os.Getenv("GRPC_SERVER_ADDR"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error listening to the grpc server : %v", err)
	}
}
