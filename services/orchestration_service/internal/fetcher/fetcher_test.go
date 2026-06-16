package fetcher

import (
	"context"
	"orchestration_service/internal/mocks/mockdb"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/server"
	"orchestration_service/internal/types"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/testify/assert"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
)

func Test_StartFetcher(t *testing.T) {

	db := repository.Db{
		Auth:    &mockdb.MockDb{DbName: enum.ServiceName_AUTH_SERVICE.String(), HasData: true},
		Tracker: &mockdb.MockDb{DbName: enum.ServiceName_TRACKER_SERVICE.String(), HasData: true},
	}

	logger := logger.NewLogger()

	targetServices := []string{enum.ServiceName_AUTH_SERVICE.String(), enum.ServiceName_TRACKER_SERVICE.String()}

	jobsChan := make(chan types.Data, 2)

	ticker := time.NewTicker(time.Second * 3)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	server := server.Server{
		Db:                    &db,
		Logger:                logger,
		FetcherJobsChan:       jobsChan,
		FetcherTargetServices: &targetServices,
		Ticker:                ticker,
		FetcherWg:             &wg,
		Ctx:                   ctx,
		CtxCancel:             cancel,
	}

	go StartFetcher(server)

	time.Sleep(time.Second * 5)

	close(jobsChan)

	for data := range jobsChan {
		assert.NoError(t, data.Err)
		assert.NotZero(t, data.DbId)
		assert.NotEmpty(t, data.Payload)
		assert.NotZero(t, data.TargetService)
		assert.NotZero(t, data.CreatedBy)
		assert.Zero(t, data.NumberOfTries)
		assert.NotZero(t, data.ServiceName)
		assert.NotZero(t, data.Task)
	}

	cancel()

	wg.Wait()
}
