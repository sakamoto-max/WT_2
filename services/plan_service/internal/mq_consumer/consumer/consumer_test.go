package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"plan_service/internal/mq_consumer/server"
	"plan_service/internal/mq_consumer/types"
	"sync"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	queueMock "github.com/sakamoto-max/rabbit_mq/mock"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_StartConsumer(t *testing.T) {
	// plan queue
	// logger
	// jobs chan
	// ctx
	// cancel

	logger := logger.NewLogger()

	jobs := make(chan types.Data, 1)

	ctx, cancel := context.WithCancel(context.Background())

	planChan := make(chan amqp.Delivery, 1)
	planQueue := queueMock.QueueMock{Channel: planChan}

	var wg sync.WaitGroup
	server := server.Server{
		Ctx:        ctx,
		PlanQueue:  &planQueue,
		Logger:     logger,
		JobsChan:   jobs,
		CtxCancel:  cancel,
		ConsumerWg: &wg,
	}

	go StartConsumer(server)

	data := mqTypes.Data{
		DbId:          "1",
		SentBy:        "test",
		TaskName:      "test_task",
		Payload:       map[string]any{"test": "test"},
		TargetService: "test_service",
	}

	dataInBytes, _ := json.Marshal(data)

	planChan <- amqp.Delivery{Body: dataInBytes}

	dataReceived := <-jobs

	time.Sleep(1 * time.Second)

	assert.Equal(t, dataReceived.DbId, data.DbId)
	assert.Equal(t, dataReceived.SentBy, data.SentBy)
	assert.Equal(t, dataReceived.TaskName, data.TaskName)
	assert.Equal(t, dataReceived.TargetService, data.TargetService)
	assert.Equal(t, dataReceived.Payload, data.Payload)

	close(planChan)
	cancel()
	wg.Wait()

	close(jobs)

	fmt.Println("test completed")
}
