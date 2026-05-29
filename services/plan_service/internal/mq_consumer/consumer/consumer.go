package consumer

import (
	"encoding/json"
	// "fmt"
	"plan_service/internal/mq_consumer/types"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
	// "github.com/sakamoto-max/wt_2_proto/shared/enum"
)

type consumer struct {
	planQueue *mq.MessageQueue
	logger    *logger.MyLogger
	jobs      chan<- types.Data
}

func NewConsumer(planQueue *mq.MessageQueue, logger *logger.MyLogger, jobs chan<- types.Data) *consumer {
	return &consumer{
		planQueue: planQueue,
		logger:    logger,
		jobs:      jobs,
	}
}

func (c *consumer) Start() {

	msgs, err := c.planQueue.Consume()
	if err != nil {
		c.logger.Log.Fatalf("unable to get data from the consumer : %v", err)
	}

	for {
		msg, ok := <-msgs

		if !ok {
			// c.logger.Log.Info("consumer is stopped")
			return
		}

		data := ConvertIntoJosn(&msg.Body)

		c.logger.Log.Infow("consumer has received a task", zap.String("task name", data.TaskName), zap.String("sent by", data.SentBy))
		
		c.jobs <- types.Data{
			DbId: data.DbId,
			SentBy: data.SentBy,
			TaskName: data.TaskName,
			Payload: data.Payload,
			TargetService: data.TargetService,
		}
		c.logger.Log.Infoln("consumer has sent data to the jobs chan")
	}

}

func ConvertIntoJosn(data *[]byte) mqTypes.Data {

	var D mqTypes.Data

	_ = json.Unmarshal(*data, &D)

	return D
}
