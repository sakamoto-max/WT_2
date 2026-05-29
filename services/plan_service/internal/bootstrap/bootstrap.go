package bootstrap

import (
	"net"
	"os"
	"os/signal"
	"plan_service/internal/client"
	"plan_service/internal/database"
	"plan_service/internal/repository/cache"
	"plan_service/internal/repository"
	"plan_service/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	pb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type app struct {
	addr        string
	service     *services.Service
	logger      *logger.MyLogger
	pool        *pgxpool.Pool
	redisClient *redis.Client
	exerConn        *grpc.ClientConn
}

func NewApp(addr string) *app {

	logger := logger.NewLogger()

	pool, err := database.NewPgConn()
	if err != nil {
		logger.Log.Fatalw("failed to open postgres pool", zap.Error(err))
	}

	redisClient, err := database.NewRedisConn()
	if err != nil {
		logger.Log.Fatalw("failed to open redis client", zap.Error(err))
	}

	pg := repository.NewDb(pool)
	cache := cache.NewCache(redisClient)

	exerConn := client.NewEmptyClient().OpenConnection(os.Getenv("EXERCISE_GRPC_SERVER_ADDR"))
	exerClient := exerConn.CreateExerciseClient()

	service := services.NewService(pg, cache, exerClient)
	// handler := handler.NewHandler(service, logger)

	return &app{
		addr:    addr,
		service: service,
		logger:  logger,
		pool: pool,
		redisClient: redisClient,
		exerConn: exerConn.Conn,
	}

}

func (a *app) Run() {

	defer func(){
		err := a.redisClient.Close()
		if err != nil {
			a.logger.Log.Infow("failed to close redis client", zap.Error(err))
		}
		a.exerConn.Close()
		if err != nil {
			a.logger.Log.Infow("failed to close exer grpc conn", zap.Error(err))
		}
		a.pool.Close()
	}()

	lis, err := net.Listen("tcp", a.addr)
	if err != nil {
		a.logger.Log.Fatalf("failed to listen to tcp : %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPlanServiceServer(grpcServer, a.service)

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
