package bootstrap

import (
	"exercise_service/internal/config"
	"exercise_service/internal/database"
	"exercise_service/internal/repository"
	"exercise_service/internal/repository/cache"
	"exercise_service/internal/services"
	"net"
	"os"
	"os/signal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	pb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type app struct {
	service     *services.Service
	logger      *logger.MyLogger
	pool        *pgxpool.Pool
	redisClient *redis.Client
	config      config.Config
}

func NewApp(config config.Config) *app {

	pool := database.NewPgConn(config)
	config.Logger.Log.Infoln("connected to postgres")

	redisClient := database.NewRedisConn(config)
	config.Logger.Log.Infoln("connected to redis")

	pg := repository.NewDb(pool)
	cache := cache.NewCache(redisClient)

	service := services.NewService(pg, cache)

	return &app{
		config:      config,
		service:     service,
		logger:      config.Logger,
		redisClient: redisClient,
		pool:        pool,
	}

}

func (a *app) Run() {

	lis, err := net.Listen("tcp", a.config.Server.GrpcServerAddr)
	if err != nil {
		a.logger.Log.Fatalf("failed to listen to tcp : %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExerciseServiceServer(grpcServer, a.service)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		a.logger.Log.Infow("grpc server has started", zap.String("addr", a.config.Server.GrpcServerAddr))
		if err := grpcServer.Serve(lis); err != nil {
			sigChan <- os.Interrupt
			a.logger.Log.Fatalf("error listening to the grpc server : %v", err)
		}
	}()

	sig := <-sigChan

	a.logger.Log.Infof("shutdown signal received : %v", sig.String())

	grpcServer.GracefulStop()
	a.logger.Log.Infoln("grpc server is stopped")

	a.redisClient.Close()
	a.logger.Log.Infoln("redis client is closed")

	a.pool.Close()
	a.logger.Log.Infoln("pg connection is closed")

	a.logger.Log.Infof("server has shutdown")
}
