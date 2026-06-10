package server

import (
	"context"
	"fmt"
	"plan_service/internal/client"
	"plan_service/internal/config"
	"plan_service/internal/database"
	"plan_service/internal/mq_consumer/types"
	"plan_service/internal/repository"
	"plan_service/internal/repository/cache"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	Ctx             context.Context
	CtxCancel       context.CancelFunc
	Logger          *logger.MyLogger
	MqConn          *amqp091.Connection
	PlanQueue       mq.QueueIface
	ResQueue        mq.QueueIface
	PgPool          *pgxpool.Pool
	Db              *repository.Db
	JobsChan        chan types.Data
	SendersChan     chan mqTypes.Data
	SenderWg        *sync.WaitGroup
	WorkerWg        *sync.WaitGroup
	ExerConn        *grpc.ClientConn
	ExerClient      client.ExerClientIface
	NumberOfWorkers int
	NumberOfSenders int
	Cache           *cache.Cache
	// config          config.Config
}

func NewServer(config config.Config) Server {

	// logger := logger.NewLogger()

	// make url := "amqp://guest:guest@localhost:5672/"
	             // amqp://guest:guest@localhost:5673/

	mqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.Mq.MqUserName,
		config.Mq.MqPass,
		config.Mq.MqHostName,
		config.Mq.MqPort,
	)

	mqConn, err := mq.NewConn(mqURL)
	if err != nil {
		config.Logger.Log.Fatalw("failed to open rabbit mq conn", zap.Error(err))
	}
	planqueue := mq.NewMessageQueue(mqConn, enum.QueueName_PLAN_QUEUE.String())
	resQueue := mq.NewMessageQueue(mqConn, enum.QueueName_RESULT_QUEUE.String())
	// resQueue := mock.MockQueue{Open: false}

	pool := database.NewPgConn(config)
	config.Logger.Log.Infoln("connected to postgres")

	redisClient := database.NewRedisConn(config)
	config.Logger.Log.Infoln("connected to redis")

	cache := cache.NewCache(redisClient)

	db := repository.NewDb(pool)
	// logger.Log.Infoln("connected to db")

	// jobs chan

	jobsChan := make(chan types.Data, config.NumberOfWorkers*2)

	// sender chan
	senderChan := make(chan mqTypes.Data, config.NumberOfSenders*2)

	var workerWg sync.WaitGroup
	var senderWg sync.WaitGroup

	exerConn := client.NewConn(config.OtherServices.ExerServiceHost, config.OtherServices.ExerServiceAddr, config.Logger)
	exerClient := client.CreateExerciseClient(exerConn)
	config.Logger.Log.Infoln("connected to exercise client")

	ctx, cancel := context.WithCancel(context.Background())

	return Server{
		Ctx:             ctx,
		CtxCancel:       cancel,
		Logger:          config.Logger,
		MqConn:          mqConn,
		PlanQueue:       planqueue,
		ResQueue:        resQueue,
		PgPool:          pool,
		Db:              db,
		JobsChan:        jobsChan,
		SendersChan:     senderChan,
		SenderWg:        &senderWg,
		WorkerWg:        &workerWg,
		ExerConn:        exerConn,
		ExerClient:      exerClient,
		NumberOfWorkers: config.NumberOfWorkers,
		NumberOfSenders: config.NumberOfSenders,
		Cache:           cache,
	}
}

func (c Server) ShutDown(signal string) {
	c.Logger.Log.Infow("signal received", zap.String("signal", signal))

	c.CtxCancel()

	c.Logger.Log.Infoln("consumer has closed")

	close(c.JobsChan)

	c.WorkerWg.Wait()
	c.Logger.Log.Infoln("workers have stopped")

	close(c.SendersChan)
	c.SenderWg.Wait()
	c.Logger.Log.Infoln("senders have stopped")

	c.MqConn.Close()

	c.PgPool.Close()
	c.Logger.Log.Infoln("db connection  closed")

	c.Logger.Log.Infoln("shutdown")
}
