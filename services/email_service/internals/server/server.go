package server

import (
	"email_service/internals/database"
	"email_service/internals/repostitory"
	"email_service/internals/services"
	"email_service/internals/types"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type Server struct {
	Logger          *logger.MyLogger
	PgPool          *pgxpool.Pool
	Db              *repostitory.Db
	MqConn          *amqp091.Connection
	EmailQueue      *queue.MessageQueue
	ResQueue        *queue.MessageQueue
	Service         *services.Service
	SenderChan      chan types.Data
	JobsChan        chan types.Data
	WorkersWg       *sync.WaitGroup
	SendersWg       *sync.WaitGroup
	NumberOfSenders int
	NumberOfWorkers int
}

func NewSever() Server {
	logger := logger.NewLogger()
	defer logger.Log.Sync()
	// db
	pool := database.NewDb(os.Getenv("POSTGRES_CONN"), logger)

	db := repostitory.RegisterDb(pool, logger)

	logger.Log.Infoln("connected to the database")
	// mq
	conn := queue.NewConn()
	logger.Log.Infoln("connected to the rabbit mq")

	emailQueue := queue.NewMessageQueue(conn, enum.QueueName_EMAIL_QUEUE.String())

	resQueue := queue.NewMessageQueue(conn, enum.QueueName_RESULT_QUEUE.String())
	// service
	service := services.NewService(logger)

	numberOfSenders := 5
	numberOfWorkers := 5
	// sender chan
	senderChan := make(chan types.Data, numberOfSenders*2)
	// jobs chan
	jobsChan := make(chan types.Data, numberOfWorkers*2)

	var senderWg sync.WaitGroup

	var workerWg sync.WaitGroup

	return Server{
		Logger:          logger,
		PgPool:          pool,
		Db:              db,
		MqConn:          conn,
		EmailQueue:      emailQueue,
		ResQueue:        resQueue,
		Service:         service,
		JobsChan:        jobsChan,
		SenderChan:      senderChan,
		WorkersWg:       &workerWg,
		SendersWg:       &senderWg,
		NumberOfSenders: numberOfSenders,
		NumberOfWorkers: numberOfWorkers,
	}
}

func (s Server) Shutdown(signal string) {

	s.Logger.Log.Infow("shutdown signal received", zap.String("signal", signal))

	// close consumer -> close the emailqueue
	err := s.EmailQueue.Ch.Close()
	if err != nil {
		s.Logger.Log.Errorw("falied to close the email queue channel", zap.Error(err))
	}

	s.Logger.Log.Infoln("consumer have stopped")
	// close worker
	close(s.JobsChan)
	s.WorkersWg.Wait()

	s.Logger.Log.Infoln("workers have stopped")
	// close senders
	close(s.SenderChan)
	s.SendersWg.Wait()

	s.Logger.Log.Infoln("senders have stopped")
	// close pool
	s.PgPool.Close()

	s.Logger.Log.Infoln("database connection is closed")
	// mq conn
	s.MqConn.Close()
	s.Logger.Log.Infoln("rabbit mq connection is closed")
}
