package consumer

import (
	"email_service/internals/mocks/queuemock"
	"email_service/internals/server"
	"email_service/internals/types"
	"encoding/json"
	"sync"
	"testing"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/stretchr/testify/assert"
)

func Test_StartConsumer(t *testing.T) {

	emailChan := make(chan amqp.Delivery)

	emailQueue := queuemock.QueueMock{Down: false, Channel: emailChan}

	jobsChan := make(chan types.Data, 1)

	logger := logger.NewLogger()

	var wg sync.WaitGroup

	server := server.Server{
		Logger:     logger,
		EmailQueue: &emailQueue,
		JobsChan:   jobsChan,
		ConsumerWg: &wg,
	}

	go StartConsumer(server)

	payload := map[string]any{
		enum.QueueFields_EMAIL.String(): "test1@gmail.com",
	}

	data := types.Data{
		DbId:          "123",
		TargetService: "targer_service",
		TaskName:      "task_name",
		Payload:       payload,
		SentBy:        "sent_by",
	}

	dataInBytes, _ := json.Marshal(data)

	emailChan <- amqp.Delivery{
		Body: dataInBytes,
	}

	dataFromJobsChan := <-jobsChan

	assert.Equal(t, dataFromJobsChan.DbId, data.DbId)
	assert.Equal(t, dataFromJobsChan.Payload, data.Payload)
	assert.Equal(t, dataFromJobsChan.SentBy, data.SentBy)
	assert.Equal(t, dataFromJobsChan.TaskName, data.TaskName)
	assert.Equal(t, dataFromJobsChan.TargetService, data.TargetService)

	close(emailChan)
	wg.Wait()
	close(jobsChan)
}

// if we send data into the emailchan -> the consumer will consume it and will
// send it to jobs chan
// have to verify if the jobs chan is getting data correctly or not
