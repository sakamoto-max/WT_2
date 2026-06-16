package consumer

import (
	"orchestration_service/internal/server"
	"orchestration_service/internal/types"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/testify/assert"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sakamoto-max/rabbit_mq/mock"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
)

func Test_StartConsumer(t *testing.T) {

	jobsChan := make(chan types.Data, 1)
	logger := logger.NewLogger()

	resultQueueChan := make(chan amqp.Delivery, 1)

	resultQueue := mock.QueueMock{
		Channel: resultQueueChan,
	}

	var wg sync.WaitGroup

	server := server.Server{
		Logger:      logger,
		ResultQueue: &resultQueue,
		ConsumerWg:  &wg,
		ConsumerJobsChan:    jobsChan,
	}

	go StartConsumer(server)

	data := mqTypes.Data{
		DbId:          "123",
		TaskName:      "task",
		TargetService: "target",
		SentBy:        "sent",
		TaskStatus:    "status",
	}

	dataInBytes, _ := data.ConvertIntoBytes()

	sendingData := amqp.Delivery{
		Body: *dataInBytes,
	}

	resultQueueChan <- sendingData

	time.Sleep(time.Second)

	dataFromJobsChan := <-jobsChan

	assert.Equal(t, dataFromJobsChan.DbId, data.DbId)
	assert.Equal(t, dataFromJobsChan.TargetService, data.TargetService)
	assert.Equal(t, dataFromJobsChan.CreatedBy, data.SentBy)
	assert.Equal(t, dataFromJobsChan.Task, data.TaskName)
	assert.Equal(t, dataFromJobsChan.Status, data.TaskStatus)

	close(resultQueueChan)
	wg.Wait()
	close(jobsChan)
}
