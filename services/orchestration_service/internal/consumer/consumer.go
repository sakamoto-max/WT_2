package consumer

import (
	"encoding/json"
	server "orchestration_service/internal/server"
	"orchestration_service/internal/types"
	"sync"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

type consumer struct {
	resultQueue mq.QueueIface
	logger      *logger.MyLogger
	jobs        chan<- types.Data
}

func StartConsumer(server server.Server) {
	c := consumer{
		resultQueue: server.ResultQueue,
		logger:      server.Logger,
		jobs:        server.JobsChan,
	}

	server.ConsumerWg.Add(1)
	go c.Start(server.ConsumerWg)

	server.Logger.Log.Infoln("consumer has started")
}

func (c *consumer) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	msgsQueue, err := c.resultQueue.Consume()
	if err != nil {
		c.logger.Log.Fatalf("unable to open the consumer : %v", err)
	}

	for {
		msg, ok := <-msgsQueue
		if !ok {
			return
		}

		var data mqTypes.Data

		_ = json.Unmarshal(msg.Body, &data)

		c.logger.Log.Infow("consumer got data",
			zap.String("task name", data.TaskName),
			zap.String("sent by", data.SentBy),
			zap.String("task status", data.TaskStatus),
		)

		c.jobs <- types.Data{
			DbId:          data.DbId,
			TargetService: data.TargetService,
			CreatedBy:     data.SentBy,
			Task:          data.TaskName,
			Status:        data.TaskStatus,
			Err:           data.Err,
		}

	}
}
