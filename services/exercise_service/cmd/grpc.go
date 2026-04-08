package main

import (
	"exercise_service/internal/controllers"
	"exercise_service/internal/repository"
	"exercise_service/internal/services"
	"log"
	"net"
	"os"
	"os/signal"
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
	// logger := logger.NewLogger()

	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		log.Fatalf("failed to listen to tcp : %v", err)
	}

	log.Printf("created TCP listener at %v", g.addr)

	repo, err := repository.NewRepo()
	if err != nil {
		lis.Close()
		log.Fatalf("error opening the repos : %v", err)
	}

	log.Println("created Db connections")

	defer func() {
		if err := repo.Close(); err != nil {
			log.Fatalf("error closing the databases : %v", err)
		}
	}()

	service := services.NewService(repo)

	controller := controllers.NewExerController(service)

	grpcServer := grpc.NewServer()
	exercisepb.RegisterExerciseServiceServer(grpcServer, controller)
	
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	go func() {
		log.Printf("grpc server has started at %v", g.addr)
		if err := grpcServer.Serve(lis); err != nil {
			sigChan <- os.Interrupt
			log.Fatalf("error listening to the grpc server : %v", err)
		}
	}()

	sig := <-sigChan

	log.Printf("shutdown signal received : %v", sig.String())

	grpcServer.GracefulStop()

	log.Println("gracefully shutdown")

}
