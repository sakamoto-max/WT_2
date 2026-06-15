package consumer

import (
	"context"
	"encoding/json"
	"plan_service/internal/mq_consumer/server"
	"plan_service/internal/mq_consumer/types"
	"sync"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

type consumer struct {
	planQueue mq.QueueIface
	logger    *logger.MyLogger
	jobs      chan<- types.Data
}

func StartConsumer(server server.Server) {

	consumer := &consumer{
		planQueue: server.PlanQueue,
		logger:    server.Logger,
		jobs:      server.JobsChan,
	}

	server.ConsumerWg.Add(1)
	go consumer.Start(server.Ctx, server.ConsumerWg)

	server.Logger.Log.Infoln("consumer has started")

}

func (c *consumer) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	msgs, err := c.planQueue.Consume()
	if err != nil {
		c.logger.Log.Fatalf("unable to get data from the consumer : %v", err)
	}

	for {

		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:

			if !ok {
				return
			}

			data := ConvertIntoJosn(&msg.Body)

			c.logger.Log.Infow("consumer has received a task",
				zap.String("task name", data.TaskName),
				zap.String("sent by", data.SentBy),
			)

			c.jobs <- types.ToData(data)

			c.logger.Log.Infoln("consumer sent data to the jobs chan")

		}

	}

}

func ConvertIntoJosn(data *[]byte) mqTypes.Data {

	var D mqTypes.Data

	_ = json.Unmarshal(*data, &D)

	return D
}
