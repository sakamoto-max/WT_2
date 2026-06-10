package bootstrap

import (
	"auth_service/internal/config"
	"auth_service/internal/database"
	"auth_service/internal/repository"
	"auth_service/internal/repository/cache"
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
	service     *services.Service
	Logger      *logger.MyLogger
	pool        *pgxpool.Pool
	redisClient *redis.Client
	config      config.Config
}

func NewApp(config config.Config) *app {

	logger := logger.NewLogger()

	pool := database.NewPgConn(config)

	logger.Log.Infoln("connected to postgres")

	redisClient := database.NewRedisConn(config)

	logger.Log.Infoln("connected to redis")

	d := repository.RegisterDB(pool)
	cache := cache.NewCache(redisClient)

	service := services.NewService(d, cache)
	return &app{

		config:      config,
		service:     service,
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

	lis, err := net.Listen("tcp", a.config.Server.GrpcServerAddr)
	if err != nil {
		a.Logger.Log.Fatalf("failed to listen to tcp : %v", err)
	}

	a.Logger.Log.Infoln("started tcp connection")

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, a.service)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		a.Logger.Log.Infow("grpc server has started", zap.String("addr", a.config.Server.GrpcServerAddr))
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
