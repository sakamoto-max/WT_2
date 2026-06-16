package workers

import (
	"encoding/json"
	"fmt"
	"orchestration_service/internal/mocks/cachemock"
	"orchestration_service/internal/mocks/mockdb"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/server"
	"orchestration_service/internal/types"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/testify/assert"
	amqp "github.com/rabbitmq/amqp091-go"
	queuemock "github.com/sakamoto-max/rabbit_mq/mock"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
)

func Test_StartWorkers(t *testing.T) {
	logger := logger.NewLogger()

	FetcherjobsChan := make(chan types.Data, 2)
	ConsumerJobsChan := make(chan types.Data, 2)

	update := make([]string, 1)

	repo := repository.Db{
		Auth:    &mockdb.MockDb{Update: update},
		Tracker: &mockdb.MockDb{},
		Cache:   &cachemock.CacheMock{},
	}

	planQueueChan := make(chan amqp.Delivery, 1)

	planQueue := queuemock.QueueMock{
		Channel: planQueueChan,
	}

	emailQueueChan := make(chan amqp.Delivery, 1)

	emailQueue := queuemock.QueueMock{
		Channel: emailQueueChan,
	}

	var fetcherWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	server := server.Server{
		PlanQueue:        &planQueue,
		EmailQueue:       &emailQueue,
		Db:               &repo,
		NumberOfWorkers:  1,
		FetcherJobsChan:  FetcherjobsChan,
		ConsumerJobsChan: ConsumerJobsChan,
		Logger:           logger,
		FetcherWorkersWg: &fetcherWg,
		ConsumerWorkerWg: &consumerWg,
	}

	go StartWorkersForFetcher(server)
	go StartWorkersForConsumer(server)

	dataOne := types.Data{
		TaskId:        "123",
		DbId:          "123",
		TargetService: enum.ServiceName_PLAN_SERVICE.String(),
		CreatedBy:     enum.ServiceName_AUTH_SERVICE.String(),
		Task:          enum.TaskName_CREATE_EMPTY_PLAN_FOR_USER.String(),
		Payload:       map[string]any{enum.QueueFields_USER_ID.String(): "123"},
	}

	dataTwo := types.Data{
		TaskId:        "456",
		DbId:          "123",
		TargetService: enum.ServiceName_EMAIL_SERVICE.String(),
		CreatedBy:     enum.ServiceName_AUTH_SERVICE.String(),
		Task:          enum.TaskName_SEND_EMAIL_FOR_SIGNING_UP.String(),
		Payload:       map[string]any{enum.QueueFields_EMAIL.String(): "test1@gmail.com"},
	}

	FetcherjobsChan <- dataOne
	FetcherjobsChan <- dataTwo

	dataThree := types.Data{
		TaskId:        "456",
		DbId:          "123",
		TargetService: enum.ServiceName_AUTH_SERVICE.String(),
		CreatedBy:     enum.ServiceName_EMAIL_SERVICE.String(),
		Task:          enum.TaskName_UPDATE_VALUE_IN_DB.String(),
		Status:        enum.TaskStatus_TASK_COMPLETED.String(),
	}

	ConsumerJobsChan <- dataThree

	time.Sleep(time.Second)

	assert.NotNil(t, planQueue.Data)
	assert.NotNil(t, emailQueue.Data)
	assert.Equal(t, enum.TaskStatus_TASK_COMPLETED.String(), update[0])

	close(FetcherjobsChan)
	server.FetcherWorkersWg.Wait()

	close(ConsumerJobsChan)
	server.ConsumerWorkerWg.Wait()

	close(planQueueChan)
	close(emailQueueChan)
}

func Test_NumberOfTriesExceeded(t *testing.T) {

	logger := logger.NewLogger()

	jobsChan := make(chan types.Data, 2)
	consumerJobsChan := make(chan types.Data, 2)

	update := make([]string, 1)
	tracker := make([]string, 1)

	repo := repository.Db{
		Auth:    &mockdb.MockDb{Update: update},
		Tracker: &mockdb.MockDb{Update: tracker},
		Cache:   &cachemock.CacheMock{},
	}

	planQueue := queuemock.QueueMock{}

	emailQueue := queuemock.QueueMock{}

	deadLetterQueue := queuemock.QueueMock{}

	var wg sync.WaitGroup
	var consumerWg sync.WaitGroup

	server := server.Server{
		PlanQueue:        &planQueue,
		EmailQueue:       &emailQueue,
		DeadLetterQueue:  &deadLetterQueue,
		Db:               &repo,
		NumberOfWorkers:  1,
		FetcherJobsChan:  jobsChan,
		Logger:           logger,
		FetcherWorkersWg: &wg,
		ConsumerJobsChan: consumerJobsChan,
		ConsumerWorkerWg: &consumerWg,
	}

	go StartWorkersForFetcher(server)
	go StartWorkersForConsumer(server)

	dataOne := types.Data{
		TaskId:        "123",
		DbId:          "123",
		TargetService: enum.ServiceName_PLAN_SERVICE.String(),
		CreatedBy:     enum.ServiceName_AUTH_SERVICE.String(),
		Task:          enum.TaskName_CREATE_EMPTY_PLAN_FOR_USER.String(),
		Payload:       map[string]any{enum.QueueFields_USER_ID.String(): "123"},
		WorkersTries:  6,
	}

	jobsChan <- dataOne
	time.Sleep(time.Second)

	d := deadLetterQueue.Data

	var dataFromQueue types.Data

	err := json.Unmarshal(*d, &dataFromQueue)
	if err != nil {
		t.Fatalf("failed to unmarshal : %s", err)
	}

	assert.Equal(t, dataFromQueue, dataOne)

	dataThree := types.Data{
		TaskId:        "456",
		DbId:          "123",
		TargetService: enum.ServiceName_AUTH_SERVICE.String(),
		CreatedBy:     enum.ServiceName_EMAIL_SERVICE.String(),
		Task:          enum.TaskName_UPDATE_VALUE_IN_DB.String(),
		Status:        enum.TaskStatus_TASK_COMPLETED.String(),
		WorkersTries:  6,
	}

	consumerJobsChan <- dataThree

	time.Sleep(time.Second * 2)

	d = deadLetterQueue.Data

	err = json.Unmarshal(*d, &dataFromQueue)
	if err != nil {
		t.Fatalf("failed to unmarshal : %s", err)
	}

	assert.Equal(t, dataFromQueue, dataThree)

	close(jobsChan)
	wg.Wait()

	close(consumerJobsChan)
	consumerWg.Wait()

	fmt.Println("test ended")
}

func Test_ExponentialBackOff(t *testing.T) {

	logger := logger.NewLogger()

	jobsChan := make(chan types.Data, 1)
	consumerJobsChan := make(chan types.Data, 1)

	cache := make([]string, 1)

	repo := repository.Db{
		Auth:    &mockdb.MockDb{},
		Tracker: &mockdb.MockDb{},
		Cache:   &cachemock.CacheMock{Skip: true, Data: cache},
	}

	planQueue := queuemock.QueueMock{}

	emailQueue := queuemock.QueueMock{}

	var wg sync.WaitGroup
	var consumerWg sync.WaitGroup

	server := server.Server{
		PlanQueue:        &planQueue,
		EmailQueue:       &emailQueue,
		Db:               &repo,
		NumberOfWorkers:  1,
		FetcherJobsChan:  jobsChan,
		Logger:           logger,
		FetcherWorkersWg: &wg,
		ConsumerJobsChan: consumerJobsChan,
		ConsumerWorkerWg: &consumerWg,
	}

	go StartWorkersForFetcher(server)
	go StartWorkersForConsumer(server)

	dataOne := types.Data{
		TaskId:        "456",
		DbId:          "123",
		TargetService: enum.ServiceName_AUTH_SERVICE.String(),
		CreatedBy:     enum.ServiceName_EMAIL_SERVICE.String(),
		Task:          enum.TaskName_UPDATE_VALUE_IN_DB.String(),
		Status:        enum.TaskStatus_TASK_COMPLETED.String(),
		Payload:       map[string]any{enum.QueueFields_USER_ID.String(): "123"},
	}

	jobsChan <- dataOne

	time.Sleep(time.Second * 2)
	
	dataFromJobs := <-jobsChan

	assert.Equal(t, dataFromJobs, dataOne)

	dataThree := types.Data{
		TaskId:        "456",
		DbId:          "123",
		TargetService: enum.ServiceName_AUTH_SERVICE.String(),
		CreatedBy:     enum.ServiceName_EMAIL_SERVICE.String(),
		Task:          enum.TaskName_UPDATE_VALUE_IN_DB.String(),
		Status:        enum.TaskStatus_TASK_COMPLETED.String(),
		WorkersTries:  6,
	}

	consumerJobsChan <- dataThree
	time.Sleep(time.Second * 2)

	dataFromConsumerJobs := <-consumerJobsChan

	assert.Equal(t, dataFromConsumerJobs, dataThree)

	close(jobsChan)
	wg.Wait()

	close(consumerJobsChan)
	consumerWg.Wait()

	fmt.Println("test ended")
}
