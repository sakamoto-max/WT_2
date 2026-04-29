package consumer

import (
	"email_service/internals/types"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type consumer struct {
	emailQueue *queue.MessageQueue
	logger     *logger.MyLogger
	jobs       chan<- types.Data
}

func NewConsumer(emailQueue *queue.MessageQueue, logger *logger.MyLogger, jobs chan<- types.Data) *consumer {
	return &consumer{
		emailQueue: emailQueue,
		logger:     logger,
		jobs:       jobs,
	}
}

func (c *consumer) Start() <-chan amqp091.Delivery {
	c.logger.Log.Infoln("consumer has started")
	msgs, err := c.emailQueue.Consume(enum.QueueName_EMAIL_QUEUE.String())
	if err != nil {
		c.logger.Log.Fatalw("failed to consume from email queue", zap.Error(err))
	}

	return msgs
}

func (c *consumer) PushToJobs(msgs <-chan amqp091.Delivery) {

	for {
		msg, ok := <-msgs

		if !ok {
			c.logger.Log.Infoln("closing the consumer")
			close(c.jobs)
			return
		}

		data := convertIntoStruct(&msg.Body)

		c.jobs <- *data
	}
}

func convertIntoStruct(data *[]byte) *types.Data {

	var D types.Data

	_ = json.Unmarshal(*data, &D)

	return &D
}
