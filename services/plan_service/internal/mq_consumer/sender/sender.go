package sender

import (
	"context"
	"plan_service/internal/config"
	"sync"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

type Sender struct {
	id         int
	senderChan <-chan mqTypes.Data
	resQueue   *mq.MessageQueue
	logger     *logger.MyLogger
}

func StartSenders(config config.Config) {

	for i := range config.NumberOfSenders {
		sender := &Sender{
			id:         i + 1,
			senderChan: config.SendersChan,
			resQueue:   config.ResQueue,
			logger:     config.Logger,
		}
		config.SenderWg.Add(1)
		go sender.Start(config.SenderWg)
	}

	config.Logger.Log.Infow("senders have started", zap.Int("number of senders", config.NumberOfSenders))
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
			s.logger.Log.Errorw("failed to send data to the result queue", zap.Int("sender Id", s.id), zap.Error(err))
			continue
		}

		s.logger.Log.Infow("sender has successfully published the data to the result queue", zap.Int("sender Id", s.id))
	}
}

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
