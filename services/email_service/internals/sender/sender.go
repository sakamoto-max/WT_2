package sender

import (
	"context"
	"email_service/internals/repostitory"
	"email_service/internals/types"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type Sender struct {
	id         int
	logger     *logger.MyLogger
	resQueue   *queue.MessageQueue
	senderChan <-chan types.Data
	db         *repostitory.Db
}

func MakeSenders(numberOfSenders int, logger *logger.MyLogger, resQueue *queue.MessageQueue, senderChan <-chan types.Data, db *repostitory.Db) []*Sender {

	var senders []*Sender

	for i := range numberOfSenders {
		s := &Sender{
			id:         i + 1,
			logger:     logger,
			resQueue:   resQueue,
			senderChan: senderChan,
			db:         db,
		}

		senders = append(senders, s)
	}

	return senders
}

func (s *Sender) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg, ok := <-s.senderChan
		if !ok {
			return
		}

		dataInBytes, _ := convertIntoBytes(msg)

		err := s.resQueue.Publish(context.Background(), dataInBytes)
		if err != nil {
			s.logger.Log.Infow("publishing to the result Queue failed. entering exponential backoff", zap.Int("sender id", s.id))
			s.ExponentialBackOff(msg)
			continue
		}

		s.logger.Log.Infow("sender have successfully pushed the data into result Queue", zap.Int("sender Id", s.id))

	}

}

func (s *Sender) ExponentialBackOff(data types.Data) {

	dataInBytes, _ := convertIntoBytes(data)

	waitTime := time.Millisecond * 300
	numberOfTries := 3

	var err error

	for range numberOfTries {

		err = s.resQueue.Publish(context.Background(), dataInBytes)
		if err != nil {
			time.Sleep(waitTime)
			waitTime = waitTime * 2
			continue
		}

		return
	}

	s.logger.Log.Infoln("exponential backoff failed. pushing data into the db", zap.Int("sender id", s.id))

	err = s.db.PushToFailed(data, numberOfTries, enum.TaskStatus_TASK_FAILED.String(), err)
	if err != nil {
		s.logger.Log.Errorw("db operation failed", zap.Error(err))
	}

}


func convertIntoBytes(data any) (*[]byte, error) {
	dataInBytes, err := json.Marshal(data)

	if err != nil {
		return nil, fmt.Errorf("failed to convert into bytes %w", err)
	}

	return &dataInBytes, nil
}
