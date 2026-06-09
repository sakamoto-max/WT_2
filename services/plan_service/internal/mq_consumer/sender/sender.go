package sender

import (
	"context"
	// "fmt"
	"plan_service/internal/mq_consumer/server"
	"plan_service/internal/repository"
	"sync"
	"time"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

type Sender struct {
	id         int
	senderChan <-chan mqTypes.Data
	resQueue   mq.QueueIface
	logger     *logger.MyLogger
	db         *repository.Db
}

func StartSenders(server server.Server) {

	for i := range server.NumberOfSenders {
		sender := &Sender{
			id:         i + 1,
			senderChan: server.SendersChan,
			resQueue:   server.ResQueue,
			logger:     server.Logger,
			db:         server.Db,
		}
		server.SenderWg.Add(1)
		go sender.Start(server.SenderWg)
	}

	server.Logger.Log.Infow("senders have started", zap.Int("number of senders", server.NumberOfSenders))
}

func (s *Sender) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		msg, ok := <-s.senderChan
		if !ok {
			return
		}

		dataInBytes, _ := msg.ConvertIntoBytes()

		err := s.resQueue.Publish(context.Background(), dataInBytes)
		if err != nil {
			s.exponentialBackOff(context.Background(), msg)
			continue
		}

		s.logger.Log.Infow("sender has successfully published the data to the result queue",
			zap.Int("sender Id", s.id),
		)
	}
}

func (s *Sender) exponentialBackOff(ctx context.Context, data mqTypes.Data) {

	waitTime := time.Millisecond * 300
	numberOfTries := 3

	dataInBytes, _ := data.ConvertIntoBytes()

	for range numberOfTries {
		err := s.resQueue.Publish(ctx, dataInBytes)
		if err != nil {
			time.Sleep(waitTime)
			waitTime = waitTime * 2
			continue
		}

		s.logger.Log.Infow("sender has successfully published the data to the result queue",
			zap.Int("sender Id", s.id),
		)

		return
	}

	if err := s.db.QueueDb.Insert(data); err != nil {
		s.logger.Log.Errorw("failed to push data to the result queue and insert into the failed_to_push_to_queue db is also failed",
			zap.Int("sender id", s.id),
			zap.Error(err),
		)

		return
	}

	s.logger.Log.Errorw("failed to push data to the result queue",
		zap.Int("sender id", s.id),
	)
}

// failed to push to queue -> store data in db
// id | targer_service| target_db_id | taskName | task_status | reason

// consumer :
// listens to the queue
// if any msg is found -> sends it to the worker

// worker :
// receives work from the consumer
// does the work
// sends data to the sender chan whether the work is passed or failed

// sender :
// listens to the sender chan
// if any work is received -> publishes it to the resultQueue
