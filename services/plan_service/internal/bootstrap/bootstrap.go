package bootstrap

import (
	"net"
	"os"
	"os/signal"
	"plan_service/internal/client"
	"plan_service/internal/config"
	"plan_service/internal/database"
	"plan_service/internal/repository"
	"plan_service/internal/repository/cache"
	"plan_service/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	pb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type app struct {
	service     *services.Service
	logger      *logger.MyLogger
	pool        *pgxpool.Pool
	redisClient *redis.Client
	exerConn    *grpc.ClientConn
	config      config.Config
}

func NewApp(config config.Config) *app {

	logger := logger.NewLogger()

	pool := database.NewPgConn(config)
	config.Logger.Log.Infoln("connected to postgres")

	redisClient := database.NewRedisConn(config)
	config.Logger.Log.Infoln("connected to redis")

	pg := repository.NewDb(pool)
	cache := cache.NewCache(redisClient)

	exerConn := client.NewConn(config.OtherServices.ExerServiceHost, config.OtherServices.ExerServiceAddr, config.Logger)
	exerClient := client.CreateExerciseClient(exerConn)
	config.Logger.Log.Infoln("connected to exercise client")

	service := services.NewService(pg, cache, exerClient)

	return &app{
		service:     service,
		logger:      logger,
		pool:        pool,
		redisClient: redisClient,
		exerConn:    exerConn,
		config:      config,
	}

}

func (a *app) Run() {

	lis, err := net.Listen("tcp", a.config.Server.GrpcServerAddr)
	if err != nil {
		a.logger.Log.Fatalf("failed to listen to tcp : %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPlanServiceServer(grpcServer, a.service)

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
	a.logger.Log.Infoln("grpc server is closed")

	a.redisClient.Close()
	a.logger.Log.Infoln("redis client is closed")

	a.exerConn.Close()
	a.logger.Log.Infoln("exercise service connection is closed")

	a.pool.Close()
	a.logger.Log.Infoln("postgres connection is closed")

	a.logger.Log.Infof("server has shutdown")
}
