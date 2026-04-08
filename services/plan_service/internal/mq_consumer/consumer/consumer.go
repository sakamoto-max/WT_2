package consumer

import (
	"wt/pkg/enum"
	"wt/pkg/logger"
	mq "wt/pkg/queue"
	"wt/pkg/types"
	"wt/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
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
	msgs, err := c.planQueue.Consume(string(enum.PlanQueue))
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

		data := utils.ConvertIntoJosn(&msg.Body)

		c.jobs <- *data

	}
}
