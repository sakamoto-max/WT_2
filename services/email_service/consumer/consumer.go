package consumer

import (

	"github.com/sakamoto-max/wt_2-pkg/enum"
	"github.com/sakamoto-max/wt_2-pkg/logger"
	// "github.com/sakamoto-max/wt_2-pkg/queue"
	"github.com/sakamoto-max/rabbit_mq/queue" 
	"github.com/sakamoto-max/wt_2-pkg/types"
	"github.com/sakamoto-max/wt_2-pkg/utils"

	"github.com/rabbitmq/amqp091-go"
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
	msgs, err := c.emailQueue.Consume(string(enum.EmailQueue))
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

		data := utils.ConvertIntoJosn(&msg.Body)

		c.jobs <- *data
	}
}
