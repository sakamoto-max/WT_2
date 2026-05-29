package fetcher

import (
	"context"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/types"
	"sync"
	"time"

	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type fetcher struct {
	db             *repository.Db
	logger         *logger.MyLogger
	targetServices *[]string
	jobsChan       chan<- types.Data
	TickerChan     <-chan time.Time
}

func NewFetcher(db *repository.Db, logger *logger.MyLogger, targetServices *[]string, jobsChan chan<- types.Data, tickerChan <-chan time.Time) *fetcher {
	return &fetcher{
		db:             db,
		logger:         logger,
		targetServices: targetServices,
		jobsChan:       jobsChan,
		TickerChan:     tickerChan,
	}
}

func (p *fetcher) Start(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	p.logger.Log.Infoln("producer has started")

	for {
		select {
		case <-p.TickerChan:
			DataChan := make(chan *[]types.Data, len(*p.targetServices))

			var fetchWg sync.WaitGroup

			for _, targetService := range *p.targetServices {
				fetchWg.Add(1)

				switch targetService {
				case enum.ServiceName_AUTH_SERVICE.String():

					go p.db.Auth.FetchData(ctx, &fetchWg, DataChan)

				case enum.ServiceName_TRACKER_SERVICE.String():

					go p.db.Tracker.FetchData(ctx, &fetchWg, DataChan)

				}
			}

			fetchWg.Wait()

			for range *p.targetServices {
				data := <-DataChan

				for _, v := range *data {
					if v.NoData {
						p.logger.Log.Infow("no rows found", zap.String("service name", v.ServiceName))
						continue
					}

					if v.Err != nil {
						p.logger.Log.Errorw("failed to fetch data", zap.String("service name", v.ServiceName), zap.Error(v.Err))
						continue
					}

					p.jobsChan <- v
				}
			}
			close(DataChan)

		case <-ctx.Done():
			return
		}
	}
}
