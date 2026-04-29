package bootstrap

import (
	"net"
	"os"
	"os/signal"
	"tracker_service/internal/client"
	"tracker_service/internal/handler"
	"tracker_service/internal/repository"
	"tracker_service/internal/services"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	pb "github.com/sakamoto-max/wt_2_proto/shared/tracker"

	"google.golang.org/grpc"
)

type app struct {
	addr    string
	handler *handler.Handler
	logger  *logger.MyLogger
}

func NewApp(addr string) *app {

	logger := logger.NewLogger()

	repo, err := repository.NewRepo()
	if err != nil {
		logger.Log.Fatalf("error opening the repos : %v", err)
	}

	clients := client.New()

	service := services.NewService(repo, clients.PlanClient, clients.ExerClient)
	handler := handler.NewHandler(service, logger)

	return &app{
		addr:    addr,
		handler: handler,
		logger:  logger,
	}

}

func (a *app) Run() {
	lis, err := net.Listen("tcp", a.addr)
	if err != nil {
		a.logger.Log.Fatalf("failed to listen to tcp : %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTrackerServiceServer(grpcServer, a.handler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		a.logger.Log.Infof("grpc server has started at %v", a.addr)
		if err := grpcServer.Serve(lis); err != nil {
			sigChan <- os.Interrupt
			a.logger.Log.Fatalf("error listening to the grpc server : %v", err)
		}
	}()

	sig := <-sigChan

	a.logger.Log.Infof("shutdown signal received : %v", sig.String())

	grpcServer.GracefulStop()

	a.logger.Log.Infof("server is closed")
}
