package worker

import (
	"email_service/internals/server"
	"email_service/internals/services"
	"email_service/internals/types"
	"sync"
	"testing"
	"time"

	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/stretchr/testify/assert"
)

func Test_StartWorkers(t *testing.T) {
	logger := logger.NewLogger()

	service := services.NewService(logger)

	JobsChan := make(chan types.Data, 1)

	senderChan := make(chan types.Data, 1)

	var wg sync.WaitGroup

	server := server.Server{
		JobsChan:   JobsChan,
		Service:    service,
		Logger:     logger,
		SenderChan: senderChan,
		WorkersWg:  &wg,
		NumberOfWorkers: 1 ,
	}

	go StartWorkers(server)

	payload := map[string]any{
		enum.QueueFields_EMAIL.String(): "test1@gmail.com",
	}

	data := types.Data{
		DbId:          "123",
		TargetService: "targer_service",
		TaskName:      enum.TaskName_SEND_EMAIL_FOR_SIGNING_UP.String(),
		Payload:       payload,
		SentBy:        "sent_by",
	}

	JobsChan <- data

	time.Sleep(time.Second * 1)

	dataFromSenderChan := <-senderChan

	assert.Equal(t, dataFromSenderChan.DbId, data.DbId)
	assert.Equal(t, dataFromSenderChan.TargetService, data.SentBy)
	assert.Equal(t, dataFromSenderChan.SentBy, data.TargetService)
	assert.Equal(t, dataFromSenderChan.Status, enum.TaskStatus_TASK_COMPLETED.String())

	close(JobsChan)
	wg.Wait()
	close(senderChan)
}
