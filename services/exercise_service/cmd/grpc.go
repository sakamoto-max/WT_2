package main

import (
	"log"
	"net"
	"os"

	"exercise_service/internal/handlers"
	"exercise_service/internal/repository"
	"exercise_service/internal/services"
	exercisepb "workout-tracker/proto/shared/exercise"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

func (g *grpcServer) Run(pool *pgxpool.Pool, client *redis.Client) {

	lis, err := net.Listen("tcp", os.Getenv("GRPC_SERVER_ADDR"))

	if err != nil {
		log.Fatalf("failed to listen to tcp : %v", err)
	}

	grpcServer := grpc.NewServer()

	repo := repository.NewRepo(pool, client)
	service := services.NewExerciseService(repo)
	handler := handlers.NewExerciseHandler(service)

	exercisepb.RegisterExerciseServiceServer(grpcServer, handler)

	log.Printf("grpc server has started at %v", os.Getenv("GRPC_SERVER_ADDR"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error listening to the grpc server : %v", err)
	}
}
