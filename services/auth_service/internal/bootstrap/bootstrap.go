package bootstrap

import (
	"auth_service/internal/database"
	"auth_service/internal/handler"
	"auth_service/internal/repository"
	"auth_service/internal/services"
	"net"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	pb "github.com/sakamoto-max/wt_2_proto/shared/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type app struct {
	Addr        string
	Handler     *handler.Handler
	Logger      *logger.MyLogger
	pool        *pgxpool.Pool
	redisClient *redis.Client
}

func NewApp(addr string) *app {

	logger := logger.NewLogger()

	pool, err := database.NewPgConn()
	if err != nil {
		logger.Log.Fatalw("failed to connet to postgres db", zap.Error(err))
	}

	redisClient, err := database.NewRedisConn()
	if err != nil {
		logger.Log.Fatalw("failed to make redis client", zap.Error(err))
	}

	repo := repository.NewRepo(pool, redisClient)

	service := services.NewService(repo)
	handler := handler.NewHandler(service, logger)

	return &app{
		Addr:        addr,
		Handler:     handler,
		Logger:      logger,
		redisClient: redisClient,
		pool:        pool,
	}

}

func (a *app) Run() {

	defer func() {
		err := a.redisClient.Close()
		if err != nil {
			a.Logger.Log.Errorw("error in closing redis", zap.Error(err))
		}
	}()

	defer a.pool.Close()

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
