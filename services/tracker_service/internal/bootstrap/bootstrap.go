package bootstrap

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"tracker_service/internal/client"
	"tracker_service/internal/config"
	"tracker_service/internal/database"
	"tracker_service/internal/repository"
	"tracker_service/internal/repository/cache"
	"tracker_service/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	pb "github.com/sakamoto-max/wt_2_proto/shared/tracker"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type app struct {
	service     *services.Service
	logger      *logger.MyLogger
	pool        *pgxpool.Pool
	redisClient *redis.Client
	planConn    *grpc.ClientConn
	exerConn    *grpc.ClientConn
	config      config.Config
}

func NewApp(config config.Config) *app {

	pool := database.NewPgConn(config)
	config.Logger.Log.Infoln("connected to postgres")

	redisClient := database.NewRedisConn(config)
	config.Logger.Log.Infoln("connected to redis")

	pg := repository.NewDb(pool)
	cache := cache.NewCache(redisClient)

	planServerUrl := fmt.Sprintf("%s:%s", config.OtherServices.PlanServiceHost, config.OtherServices.PlanServiceAddr)

	planConn := client.OpenConnection(planServerUrl)
	planClient := client.CreatePlanClient(planConn)

	config.Logger.Log.Infoln("connected to plan service")

	exerServerUrl := fmt.Sprintf("%s:%s", config.OtherServices.ExerServiceHost, config.OtherServices.ExerServiceAddr)

	exerConn := client.OpenConnection(exerServerUrl)
	exerClient := client.CreateExerciseClient(exerConn)
	config.Logger.Log.Infoln("connected to exercise service")

	service := services.NewService(pg, cache, planClient, exerClient)

	return &app{
		service:     service,
		logger:      config.Logger,
		pool:        pool,
		redisClient: redisClient,
		planConn:    planConn,
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
	pb.RegisterTrackerServiceServer(grpcServer, a.service)

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
	a.logger.Log.Infoln("grpc server has stopped")

	a.redisClient.Close()
	a.logger.Log.Infoln("redis connection is closed")

	a.pool.Close()
	a.logger.Log.Infoln("postgres connection is closed")

	a.exerConn.Close()
	a.logger.Log.Infoln("exercise service connection is closed")

	a.planConn.Close()
	a.logger.Log.Infoln("plan service connection is closed")

	a.logger.Log.Infof("server has shutdown")
}
