package server

import (
	"context"
	"fmt"
	"orchestration_service/internal/config"
	"orchestration_service/internal/database"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/types"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sakamoto-max/rabbit_mq/queue"
	mq "github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type Server struct {
	Logger                *logger.MyLogger
	AuthPool              *pgxpool.Pool
	TrackerPool           *pgxpool.Pool
	RedisClient           *redis.Client
	Db                    *repository.Db
	MqConn                *amqp091.Connection
	PlanQueue             queue.QueueIface
	EmailQueue            queue.QueueIface
	ResultQueue           queue.QueueIface
	DeadLetterQueue       queue.QueueIface
	FetcherJobsChan       chan types.Data
	ConsumerJobsChan      chan types.Data
	Ctx                   context.Context
	CtxCancel             context.CancelFunc
	ConsumerWg            *sync.WaitGroup
	FetcherWorkersWg      *sync.WaitGroup
	ConsumerWorkerWg      *sync.WaitGroup
	FetcherWg             *sync.WaitGroup
	Ticker                *time.Ticker
	NumberOfWorkers       int
	FetcherTargetServices *[]string
}

func NewServer(config config.Config) Server {

	logger := logger.NewLogger()

	authURl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.Dbs.Auth.PgUser,
		config.Dbs.Auth.PgPass,
		config.Dbs.Auth.PgHost,
		config.Dbs.Auth.PgPort,
		config.Dbs.Auth.PgDatabaseName,
		config.Dbs.Auth.PgSSLMode,
	)

	authPool := database.NewPgConn(authURl, config)

	trackerURl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.Dbs.Tracker.PgUser,
		config.Dbs.Tracker.PgPass,
		config.Dbs.Tracker.PgHost,
		config.Dbs.Tracker.PgPort,
		config.Dbs.Tracker.PgDatabaseName,
		config.Dbs.Tracker.PgSSLMode,
	)

	trackerPool := database.NewPgConn(trackerURl, config)

	config.Logger.Log.Infoln("connected to postgres")

	redisClient, err := database.NewRedisConn(config)
	if err != nil {
		config.Logger.Log.Fatalw("failed to connect to redis", zap.Error(err))
	}

	mqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.Mq.MqUserName,
		config.Mq.MqPass,
		config.Mq.MqHostName,
		config.Mq.MqPort,
	)

	mqConn, err := mq.NewConn(mqURL)
	if err != nil {
		config.Logger.Log.Fatalw("failed to connect to rabbit mq", zap.Error(err))
	}

	planQueue := mq.NewMessageQueue(mqConn, enum.QueueName_PLAN_QUEUE.String())

	emailQueue := mq.NewMessageQueue(mqConn, enum.QueueName_EMAIL_QUEUE.String())

	resultQueue := mq.NewMessageQueue(mqConn, enum.QueueName_RESULT_QUEUE.String())

	deadLetterQueue := mq.NewMessageQueue(mqConn, "deadLetterQueue")

	authDb := repository.RegisterDb(authPool, enum.ServiceName_AUTH_SERVICE.String())
	trackerDb := repository.RegisterDb(trackerPool, enum.ServiceName_TRACKER_SERVICE.String())

	Db := repository.NewDb(authDb, trackerDb, redisClient)

	fetcherJobs := make(chan types.Data, config.Consumer.NumberOfWorkers*2)
	consumerJobs := make(chan types.Data, config.Consumer.NumberOfWorkers*2)

	ctx, cancel := context.WithCancel(context.Background())

	var fetcherWorkersWg sync.WaitGroup
	var consumerWorkersWg sync.WaitGroup
	var fetcherWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	ticker := time.NewTicker(time.Second * 30)

	targetServices := []string{enum.ServiceName_AUTH_SERVICE.String(), enum.ServiceName_TRACKER_SERVICE.String()}

	return Server{
		Logger:                logger,
		AuthPool:              authPool,
		TrackerPool:           trackerPool,
		RedisClient:           redisClient,
		Db:                    &Db,
		MqConn:                mqConn,
		PlanQueue:             planQueue,
		EmailQueue:            emailQueue,
		ResultQueue:           resultQueue,
		DeadLetterQueue:       deadLetterQueue,
		FetcherJobsChan:       fetcherJobs,
		ConsumerJobsChan:      consumerJobs,
		Ctx:                   ctx,
		CtxCancel:             cancel,
		FetcherWorkersWg:      &fetcherWorkersWg,
		ConsumerWorkerWg:      &consumerWorkersWg,
		FetcherWg:             &fetcherWg,
		ConsumerWg:            &consumerWg,
		Ticker:                ticker,
		NumberOfWorkers:       config.Consumer.NumberOfWorkers,
		FetcherTargetServices: &targetServices,
	}
}

func (s Server) Shutdown(signal string) {
	s.Logger.Log.Infow("shutdown signal received", zap.String("signal", signal))

	s.Ticker.Stop() // stops producer
	s.CtxCancel()

	s.FetcherWg.Wait()
	s.Logger.Log.Infoln("producer have stopped")

	s.MqConn.Close() // stops consumer
	s.ConsumerWg.Wait()
	s.Logger.Log.Infoln("consumer has closed")

	close(s.FetcherJobsChan)
	s.FetcherWorkersWg.Wait()
	s.Logger.Log.Infoln("fetcher workers have stopped")

	close(s.ConsumerJobsChan)
	s.ConsumerWorkerWg.Wait()
	s.Logger.Log.Infoln("consumer workers have stopped")

}
