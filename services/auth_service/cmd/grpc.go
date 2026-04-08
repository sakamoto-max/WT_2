package main

import (
	"auth_service/internal/controllers"
	"auth_service/internal/repository"
	"auth_service/internal/services"
	"net"
	"os"
	"os/signal"
	"log"
	"wt/pkg/logger"

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

	log.Print("created Db connections")

	defer func() {
		if err := repo.Close(); err != nil {
			log.Fatalf("error closing the databases : %v", err)
		}
	}()

	service := services.NewService(repo)
	logger := logger.NewLogger()
	controller := controllers.NewAuthController(service, logger)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, controller)

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

	log.Println("shutdown signal received : %v", sig.String())

	grpcServer.GracefulStop()

	log.Println("server is closed")
}
