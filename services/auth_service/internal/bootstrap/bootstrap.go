package bootstrap

import (
	"auth_service/internal/handler"
	"auth_service/internal/repository"
	"auth_service/internal/services"
	"net"
	"os"
	"os/signal"
	"wt/pkg/logger"
	pb "workout-tracker/proto/shared/auth"
	"google.golang.org/grpc"
)

type app struct {
	Addr    string
	Handler *handler.Handler
	Logger  *logger.MyLogger
}

func NewApp(addr string) *app {

	logger := logger.NewLogger()

	repo, err := repository.NewRepo()
	if err != nil {
		logger.Log.Fatalf("error opening the repos : %v", err)
	}

	service := services.NewService(repo)
	handler := handler.NewHandler(service, logger)

	return &app{
		Addr:    addr,
		Handler: handler,
		Logger:  logger,
	}

}

func (a *app) Run() {
	lis, err := net.Listen("tcp", a.Addr)
	if err != nil {
		a.Logger.Log.Fatalf("failed to listen to tcp : %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, a.Handler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		a.Logger.Log.Infof("grpc server has started at %v", a.Addr)
		if err := grpcServer.Serve(lis); err != nil {
			sigChan <- os.Interrupt
			a.Logger.Log.Fatalf("error listening to the grpc server : %v", err)
		}
	}()

	sig := <-sigChan

	a.Logger.Log.Infof("shutdown signal received : %v", sig.String())

	grpcServer.GracefulStop()

	a.Logger.Log.Infof("server is closed")
}
