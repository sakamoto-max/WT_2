package sender

import (
	"encoding/json"
	"fmt"
	"plan_service/internal/mq_consumer/server"
	"sync"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	queuemock "github.com/sakamoto-max/rabbit_mq/mock"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_StartSenders(t *testing.T) {
	// sender chan
	// res queue
	//logger
	// db
	// sender wg

	logger := logger.NewLogger()

	senderChan := make(chan mqTypes.Data, 1)

	resChan := make(chan amqp.Delivery, 1)

	resQueue := queuemock.QueueMock{
		Channel: resChan,
	}

	var senderWg sync.WaitGroup

	server := server.Server{
		NumberOfSenders: 1,
		SendersChan:     senderChan,
		ResQueue:        &resQueue,
		Logger:          logger,
		SenderWg:        &senderWg,
	}

	go StartSenders(server)

	data := mqTypes.Data{
		DbId:          "1",
		SentBy:        "test",
		TaskName:      "test_task",
		Payload:       map[string]any{"test": "test"},
		TargetService: "test_service",
	}

	senderChan <- data

	time.Sleep(time.Second * 1)

	dataGot := resQueue.Data

	var dataReceived mqTypes.Data

	json.Unmarshal(*dataGot, &dataReceived)

	assert.Equal(t, dataReceived.DbId, data.DbId)
	assert.Equal(t, dataReceived.SentBy, data.SentBy)
	assert.Equal(t, dataReceived.TaskName, data.TaskName)
	assert.Equal(t, dataReceived.TargetService, data.TargetService)
	assert.Equal(t, dataReceived.Payload, data.Payload)

	close(senderChan)
	senderWg.Wait()
	close(resChan)

	fmt.Println("test completed")
}
