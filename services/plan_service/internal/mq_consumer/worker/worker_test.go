package worker

import (
	"fmt"
	clientMock "plan_service/internal/mq_consumer/mock/mockclient"
	"plan_service/internal/mq_consumer/server"
	"plan_service/internal/mq_consumer/types"
	"plan_service/internal/repository"
	repoMock "plan_service/internal/repository/mock"
	"sync"
	"testing"
	"time"

	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/stretchr/testify/assert"
)

func Test_StartWorkers(t *testing.T) {

	logger := logger.NewLogger()

	jobsChan := make(chan types.Data, 1)

	senderChan := make(chan mqTypes.Data, 1)

	mockClient := clientMock.MockExerClient{
		ExerExists: true,
	}

	repoMock := repository.Db{
		PlanCommandRepo: &repoMock.PlanCommandMock{
			PlanExists: true,
		},
		PlanExericseRepo: &repoMock.PlanExericseMock{},
		PlanQueryRepo: &repoMock.PlanQueryMock{
			PlanExits: true,
		},
	}

	var workerWg sync.WaitGroup

	server := server.Server{
		Db:              &repoMock,
		Logger:          logger,
		JobsChan:        jobsChan,
		SendersChan:     senderChan,
		ExerClient:      &mockClient,
		NumberOfWorkers: 1,
		WorkerWg:        &workerWg,
	}

	go StartWorkers(server)

	dataToJobs := types.Data{
		DbId:          "123",
		SentBy:        enum.ServiceName_AUTH_SERVICE.String(),
		TaskName:      enum.TaskName_CREATE_EMPTY_PLAN_FOR_USER.String(),
		Payload:       map[string]any{enum.QueueFields_USER_ID.String(): "123"},
		TargetService: enum.ServiceName_PLAN_SERVICE.String(),
	}

	jobsChan <- dataToJobs

	time.Sleep(time.Second * 1)

	dataFromSenderChan := <-senderChan

	assert.Equal(t, dataToJobs.DbId, dataFromSenderChan.DbId)
	assert.Equal(t, enum.TaskName_UPDATE_VALUE_IN_DB.String(), dataFromSenderChan.TaskName)
	assert.Equal(t, dataToJobs.SentBy, dataFromSenderChan.TargetService)
	assert.Equal(t, dataToJobs.TargetService, dataFromSenderChan.SentBy)
	assert.Equal(t, dataFromSenderChan.TaskStatus, enum.TaskStatus_TASK_COMPLETED.String())

	close(jobsChan)
	workerWg.Wait()
	close(senderChan)

	fmt.Println("test completed")

}
