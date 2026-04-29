package consumer

import (
	"encoding/json"
	"fmt"
	"plan_service/internal/mq_consumer/types"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	mq "github.com/sakamoto-max/rabbit_mq/queue" 
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

func (c *consumer) GetData() <-chan amqp.Delivery {
	msgs, err := c.planQueue.Consume(enum.QueueName_PLAN_QUEUE.String())
	if err != nil {
		c.logger.Log.Fatalf("unable to get data from the consumer : %v", err)
	}

	return msgs
}

func (c *consumer) PushDataToJobs(msgs <-chan amqp.Delivery) {
	for {
		msg, ok := <-msgs

		if !ok {
			c.logger.Log.Info("closing the consumer")
			close(c.jobs)
			return
		}

		data := ConvertIntoJosn(&msg.Body)
		fmt.Println("data.Id",data.DbId)

		c.jobs <- data
	}
}

func ConvertIntoJosn(data *[]byte) types.Data {

	var D types.Data

	_ = json.Unmarshal(*data, &D)

	return D
}
