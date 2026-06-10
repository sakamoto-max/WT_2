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
	"github.com/sakamoto-max/rabbit_mq/queue"
	mq "github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type Server struct {
	Logger          *logger.MyLogger
	AuthPool        *pgxpool.Pool
	TrackerPool     *pgxpool.Pool
	Db              *repository.Db
	MqConn          *amqp091.Connection
	PlanQueue       queue.QueueIface
	EmailQueue      queue.QueueIface
	ResultQueue     queue.QueueIface
	JobsChan        chan types.Data
	Ctx             context.Context
	CtxCancel       context.CancelFunc
	WorkersWg       *sync.WaitGroup
	FetcherWg       *sync.WaitGroup
	Ticker          *time.Ticker
	NumberOfWorkers int
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

	authDb := repository.RegisterDb(authPool, enum.ServiceName_AUTH_SERVICE.String())
	trackerDb := repository.RegisterDb(trackerPool, enum.ServiceName_TRACKER_SERVICE.String())

	Db := repository.NewDb(authDb, trackerDb)

	// NumberOfWorkers := 5

	jobs := make(chan types.Data, config.Consumer.NumberOfWorkers*2)
	ctx, cancel := context.WithCancel(context.Background())

	var workerWg sync.WaitGroup
	var fetcherWg sync.WaitGroup

	ticker := time.NewTicker(time.Second * 30)

	return Server{
		Logger:          logger,
		AuthPool:        authPool,
		TrackerPool:     trackerPool,
		Db:              &Db,
		MqConn:          mqConn,
		PlanQueue:       planQueue,
		EmailQueue:      emailQueue,
		ResultQueue:     resultQueue,
		JobsChan:        jobs,
		Ctx:             ctx,
		CtxCancel:       cancel,
		WorkersWg:       &workerWg,
		FetcherWg:       &fetcherWg,
		Ticker:          ticker,
		NumberOfWorkers: config.Consumer.NumberOfWorkers,
	}
}

func (s Server) Shutdown(signal string) {
	s.Logger.Log.Infow("shutdown signal received", zap.String("signal", signal))

	s.Ticker.Stop() // stops producer
	s.CtxCancel()

	s.FetcherWg.Wait()
	s.Logger.Log.Infoln("producer have stopped")

	s.MqConn.Close() // stops consumer
	s.Logger.Log.Infoln("consumer has closed")

	close(s.JobsChan)
	s.WorkersWg.Wait()
	s.Logger.Log.Infoln("workers have stopped")
}
