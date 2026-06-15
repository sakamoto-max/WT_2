package server

import (
	"email_service/internals/config"
	"email_service/internals/database"
	"email_service/internals/repostitory"
	"email_service/internals/services"
	"email_service/internals/types"
	"fmt"
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
	Db              repostitory.RepoIFace
	MqConn          *amqp091.Connection
	EmailQueue      queue.QueueIface
	ResQueue        queue.QueueIface
	Service         *services.Service
	SenderChan      chan types.Data
	JobsChan        chan types.Data
	WorkersWg       *sync.WaitGroup
	SendersWg       *sync.WaitGroup
	ConsumerWg      *sync.WaitGroup
	NumberOfSenders int
	NumberOfWorkers int
}

func NewSever(config config.Config) Server {

	pool := database.NewPgConn(config)

	db := repostitory.RegisterDb(pool, config.Logger)

	config.Logger.Log.Infoln("connected to the database")

	mqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.Mq.UserName,
		config.Mq.Pass,
		config.Mq.Host,
		config.Mq.Port,
	)
	conn, err := queue.NewConn(mqURL)
	if err != nil {
		config.Logger.Log.Fatalw("failed to open rabbimq connection", zap.Error(err))
	}

	config.Logger.Log.Infoln("connected to the rabbit mq")

	emailQueue := queue.NewMessageQueue(conn, enum.QueueName_EMAIL_QUEUE.String())

	resQueue := queue.NewMessageQueue(conn, enum.QueueName_RESULT_QUEUE.String())
	// service
	service := services.NewService(config.Logger)

	// sender chan
	senderChan := make(chan types.Data, config.Consumer.NumberOfSenders*2)
	// jobs chan
	jobsChan := make(chan types.Data, config.Consumer.NumberOfWorkers*2)

	var senderWg sync.WaitGroup

	var workerWg sync.WaitGroup

	var consumerWg sync.WaitGroup

	return Server{
		Logger:          config.Logger,
		PgPool:          pool,
		Db:              db,
		MqConn:          conn,
		EmailQueue:      emailQueue,
		ResQueue:        resQueue,
		Service:         service,
		JobsChan:        jobsChan,
		SenderChan:      senderChan,
		ConsumerWg:      &consumerWg,
		WorkersWg:       &workerWg,
		SendersWg:       &senderWg,
		NumberOfSenders: config.Consumer.NumberOfSenders,
		NumberOfWorkers: config.Consumer.NumberOfWorkers,
	}
}

func (s Server) Shutdown(signal string) {

	s.Logger.Log.Infow("shutdown signal received", zap.String("signal", signal))

	// close consumer -> close the emailqueue
	s.EmailQueue.Close()
	s.ConsumerWg.Wait()
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
