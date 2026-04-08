package producer

import (
	"context"
	"errors"
	"orchestration_service/repository"
	"orchestration_service/types"
	"sync"
	"time"
	"wt/pkg/logger"
	mq "wt/pkg/queue"
)

type producer struct {
	db         *repository.DB
	planQueue  *mq.MessageQueue
	emailQueue *mq.MessageQueue
	logger     *logger.MyLogger
}

func NewProducer(db *repository.DB, planQueue *mq.MessageQueue, emailQueue *mq.MessageQueue, logger *logger.MyLogger) *producer {
	return &producer{
		db:         db,
		planQueue:  planQueue,
		emailQueue: emailQueue,
		logger:     logger,
	}
}

func (p *producer) Start(ctx context.Context, wg *sync.WaitGroup, tickerChan <-chan time.Time, jobsChan chan<- types.Data) {
	p.logger.Log.Infoln("producer has started")

	for {
		select {
		case <-tickerChan:
			Data, err := p.db.FetchDataFromAuth()
			if err != nil {
				if errors.Is(err, repository.ErrNoRowsFound) {
					p.logger.Log.Infoln("no rows found")
					continue
				}
				p.logger.Log.Errorln(err)
			} else {
				for i := range Data {
					jobsChan <- Data[i]
				}
			}
		case <-ctx.Done():
			wg.Wait()
			p.logger.Log.Infoln("producer stopped")
			return
		}
	}
}
