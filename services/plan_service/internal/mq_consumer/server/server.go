package server

import (
	"context"
	"os"
	"plan_service/internal/client"
	"plan_service/internal/database"
	"plan_service/internal/repository/cache"

	// "plan_service/internal/mq_consumer/mock"
	"plan_service/internal/mq_consumer/types"
	"plan_service/internal/repository"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
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
	ExerConn        *client.Client
	ExerClient      client.ExerClientIface
	NumberOfWorkers int
	NumberOfSenders int
	Cache *cache.Cache
}

func NewServer() Server {

	logger := logger.NewLogger()

	mqConn := mq.NewConn()
	planqueue := mq.NewMessageQueue(mqConn, enum.QueueName_PLAN_QUEUE.String())
	resQueue := mq.NewMessageQueue(mqConn, enum.QueueName_RESULT_QUEUE.String())
	// resQueue := mock.MockQueue{Open: false}

	pool, err := database.NewPgConn()
	if err != nil {
		logger.Log.Fatalw("failed to open postgres connection for plan consumer", zap.Error(err))
	}
	
	redisClient, err := database.NewRedisConn()
	if err != nil {
		logger.Log.Fatalw("failed to create redis connection for plan consumer", zap.Error(err))
	}

	cache := cache.NewCache(redisClient)

	db := repository.NewDb(pool)
	logger.Log.Infoln("connected to db")

	// jobs chan
	numberOfWorkers := 5
	numberOfSenders := 5

	jobsChan := make(chan types.Data, numberOfWorkers*2)

	// sender chan
	senderChan := make(chan mqTypes.Data, numberOfSenders*2)

	var workerWg sync.WaitGroup
	var senderWg sync.WaitGroup

	exerConn := client.NewEmptyClient().OpenConnection(os.Getenv("EXERCISE_GRPC_SERVER_ADDR"))
	exerciseClient := exerConn.CreateExerciseClient()

	ctx, cancel := context.WithCancel(context.Background())

	return Server{
		Ctx:             ctx,
		CtxCancel:       cancel,
		Logger:          logger,
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
		ExerClient:      exerciseClient,
		NumberOfWorkers: numberOfWorkers,
		NumberOfSenders: numberOfSenders,
		Cache: cache,
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
