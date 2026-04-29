package producer

import (
	"context"
	"errors"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/types"
	"sync"
	"time"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	mq "github.com/sakamoto-max/rabbit_mq/queue" 
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

func (p *producer) Start(ctx context.Context, wg *sync.WaitGroup, tickerChan <-chan time.Time, jobsChan chan<- types.Data, targetServices *[]string) {

	p.logger.Log.Infoln("producer has started")

	for {
		select {
		case <-tickerChan:
			p.logger.Log.Infoln("ticker signal received")

			DataChan := make(chan *[]types.Data, len(*targetServices))

			var FetchWg sync.WaitGroup

			for _, targetService := range *targetServices {
				FetchWg.Add(1)

				p.logger.Log.Infof("sent one worker to fetch data from %v", targetService)

				go FetchData(p, targetService, DataChan, &FetchWg)
			}

			FetchWg.Wait()
			p.logger.Log.Infoln("finished fetching")

			for i := 0; i < len(*targetServices); i++ {
				data := <-DataChan
				switch data {
				case nil:
					p.logger.Log.Infoln("no rows found") // in what service?
				default:
					for _, v := range *data {
						p.logger.Log.Infoln("sent data to jobs chan")
						jobsChan <- v
					}
				}
			}

			close(DataChan)

		case <-ctx.Done():
			wg.Wait()
			p.logger.Log.Infoln("producer stopped")
			return
		}
	}
}

func FetchData(p *producer, targetService string, dataChan chan<- *[]types.Data, wg *sync.WaitGroup) {

	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()

	Data, err := p.db.FetchData(ctx, targetService)

	if err != nil {
		if errors.Is(err, repository.ErrNoRowsFound) {
			dataChan <- nil
			return
		}
		p.logger.Log.Errorln(err)
		return
	}

	dataChan <- Data
}
