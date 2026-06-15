package workers

import (
	"fmt"
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

// startWorkers() -> starts the workers
//     start() -> starts listening to the jobs chan
//     pushToQueue()
// numberOfTriesExceeded()
// updateValueInDb()

// exponentialBackOff()


func Test_StartWorkers(t *testing.T) {
	logger := logger.NewLogger()

	jobsChan := make(chan types.Data, 2)

	repo := repository.Db{
		Auth:    &mockdb.MockDb{},
		Tracker: &mockdb.MockDb{},
	}

	planQueueChan := make(chan amqp.Delivery, 1)

	planQueue := queuemock.QueueMock{
		Channel: planQueueChan,
	}

	emailQueueChan := make(chan amqp.Delivery, 1)

	emailQueue := queuemock.QueueMock{
		Channel: emailQueueChan,
	}

	var wg sync.WaitGroup

	server := server.Server{
		PlanQueue:       &planQueue,
		EmailQueue:      &emailQueue,
		Db:              &repo,
		NumberOfWorkers: 1,
		JobsChan:        jobsChan,
		Logger:          logger,
		WorkersWg:       &wg,
	}

	go StartWorkers(server)

	dataOne := types.Data{
		DbId:          "123",
		TargetService: enum.ServiceName_PLAN_SERVICE.String(),
		CreatedBy:     enum.ServiceName_AUTH_SERVICE.String(),
		Task:          enum.TaskName_CREATE_EMPTY_PLAN_FOR_USER.String(),
		Payload:       map[string]any{enum.QueueFields_USER_ID.String(): "123"},
	}

	dataTwo := types.Data{
		DbId:          "123",
		TargetService: enum.ServiceName_EMAIL_SERVICE.String(),
		CreatedBy:     enum.ServiceName_AUTH_SERVICE.String(),
		Task:          enum.TaskName_SEND_EMAIL_FOR_SIGNING_UP.String(),
		Payload:       map[string]any{enum.QueueFields_EMAIL.String(): "test1@gmail.com"},
	}

	jobsChan <- dataOne
	jobsChan <- dataTwo

	time.Sleep(time.Second)

	assert.NotNil(t, planQueue.Data)
	assert.NotNil(t, emailQueue.Data)

	close(jobsChan)
	wg.Wait()

	close(planQueueChan)
	close(emailQueueChan)
}

func Test_NumberOfTriesExceeded(t *testing.T) {

	logger := logger.NewLogger()

	jobsChan := make(chan types.Data, 2)

	update := make([]string, 1)
	tracker := make([]string, 1)

	repo := repository.Db{
		Auth:    &mockdb.MockDb{Update: update},
		Tracker: &mockdb.MockDb{Update: tracker},
	}

	planQueue := queuemock.QueueMock{}

	emailQueue := queuemock.QueueMock{}

	var wg sync.WaitGroup

	server := server.Server{
		PlanQueue:       &planQueue,
		EmailQueue:      &emailQueue,
		Db:              &repo,
		NumberOfWorkers: 1,
		JobsChan:        jobsChan,
		Logger:          logger,
		WorkersWg:       &wg,
	}

	go StartWorkers(server)

	dataOne := types.Data{
		DbId:          "123",
		TargetService: enum.ServiceName_PLAN_SERVICE.String(),
		CreatedBy:     enum.ServiceName_AUTH_SERVICE.String(),
		Task:          enum.TaskName_CREATE_EMPTY_PLAN_FOR_USER.String(),
		Payload:       map[string]any{enum.QueueFields_USER_ID.String(): "123"},
		NumberOfTries: 4,
	}

	jobsChan <- dataOne

	time.Sleep(time.Second)

	assert.Equal(t, enum.TaskStatus_TASK_FAILED.String(), update[0])

	close(jobsChan)
	wg.Wait()

	fmt.Println("test ended")
}

func Test_UpdateValueInDb(t *testing.T) {

	logger := logger.NewLogger()

	jobsChan := make(chan types.Data, 2)

	auth := make([]string, 1)
	tracker := make([]string, 1)

	repo := repository.Db{
		Auth:    &mockdb.MockDb{Update: auth},
		Tracker: &mockdb.MockDb{Update: tracker},
	}

	planQueue := queuemock.QueueMock{}

	emailQueue := queuemock.QueueMock{}

	var wg sync.WaitGroup

	server := server.Server{
		PlanQueue:       &planQueue,
		EmailQueue:      &emailQueue,
		Db:              &repo,
		NumberOfWorkers: 1,
		JobsChan:        jobsChan,
		Logger:          logger,
		WorkersWg:       &wg,
	}

	go StartWorkers(server)

	dataOne := types.Data{
		DbId:          "123",
		TargetService: enum.ServiceName_AUTH_SERVICE.String(),
		CreatedBy:     enum.ServiceName_EMAIL_SERVICE.String(),
		Task:          enum.TaskName_UPDATE_VALUE_IN_DB.String(),
		Status:        enum.TaskStatus_TASK_COMPLETED.String(),
		Payload:       map[string]any{enum.QueueFields_USER_ID.String(): "123"},
	}

	jobsChan <- dataOne

	time.Sleep(time.Second)

	assert.Equal(t, enum.TaskStatus_TASK_COMPLETED.String(), auth[0])

	close(jobsChan)
	wg.Wait()

	fmt.Println("test ended")

}

func Test_ExponentialBackOff(t *testing.T) {
	// start the workers
	// disable the queue
	// make sure that the value in db is updated
}
