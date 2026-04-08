package main

import (
	"net"
	"os"
	"os/signal"
	grpcclient "plan_service/grpc_client"
	"plan_service/internal/controllers"
	"plan_service/internal/repository"
	"plan_service/internal/services"
	pb "workout-tracker/proto/shared/plan"

	"wt/pkg/logger"

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
	logger := logger.NewLogger()

	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		logger.Log.Fatalf("failed to listen to tcp : %v", err)
	}

	logger.Log.Infof("created TCP listener at %v", g.addr)

	exerClient := grpcclient.New()

	repo, err := repository.NewRepo()
	if err != nil {
		lis.Close()
		logger.Log.Fatalf("error opening the repos : %v", err)
	}

	logger.Log.Info("created Db connections")

	defer func() {
		if err := repo.Close(); err != nil {
			logger.Log.Warnf("error closing the databases : %v", err)
		}
	}()

	service := services.NewService(repo, exerClient.Client)
	controller := controllers.NewPlanController(service)

	grpcServer := grpc.NewServer()
	pb.RegisterPlanServiceServer(grpcServer, controller)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		logger.Log.Infof("grpc server has started at %v", g.addr)
		if err := grpcServer.Serve(lis); err != nil {
			sigChan <- os.Interrupt
			logger.Log.Warnf("error listening to the grpc server : %v", err)
		}
	}()

	sig := <-sigChan

	logger.Log.Infof("shutdown signal received : %v", sig.String())

	grpcServer.GracefulStop()

	logger.Log.Info("server is closed")
	
	exerClient.Close()
}
