package sender

import (
	"email_service/internals/mocks/dbmock"
	"email_service/internals/mocks/queuemock"
	"email_service/internals/server"
	"email_service/internals/types"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_StartSender(t *testing.T) {

	resultQueue := queuemock.QueueMock{Down: false}

	db := dbmock.DbMock{}

	logger := logger.NewLogger()

	senderChan := make(chan types.Data, 1)

	data := types.Data{
		DbId:          "123",
		TargetService: "targer_service",
		TaskName:      "task_name",
		SentBy:        "sent_by",
		Status:        "status",
	}

	senderChan <- data

	var wg sync.WaitGroup

	server := server.Server{
		Db:              &db,
		SenderChan:      senderChan,
		ResQueue:        &resultQueue,
		Logger:          logger,
		SendersWg:       &wg,
		NumberOfSenders: 1,
	}


	go StartSenders(server)

	var PublishedData types.Data

	time.Sleep(time.Second * 1)

	err := json.Unmarshal(*resultQueue.Data, &PublishedData)
	if err != nil {
		t.Fatalf("failed to unmarshal data")
	}

	assert.Equal(t, PublishedData.DbId, data.DbId)
	assert.Equal(t, PublishedData.TargetService, data.TargetService)
	assert.Equal(t, PublishedData.TaskName, data.TaskName)
	assert.Equal(t, PublishedData.SentBy, data.SentBy)
	assert.Equal(t, PublishedData.Status, data.Status)
	close(senderChan)
	wg.Wait()
}
