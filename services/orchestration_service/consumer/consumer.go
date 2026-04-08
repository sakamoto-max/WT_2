package consumer

import (
	"context"
	"encoding/json"
	"orchestration_service/repository"
	"wt/pkg/enum"
	"wt/pkg/logger"
	"wt/pkg/queue"
	mq "wt/pkg/queue"
	// "wt/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	db          *repository.DB
	resultQueue *mq.MessageQueue
	logger      *logger.MyLogger
}

func NewConsumer(db *repository.DB, resQueue *mq.MessageQueue, logger *logger.MyLogger) *consumer {
	return &consumer{
		db:          db,
		resultQueue: resQueue,
		logger:      logger,
	}
}


func (c *consumer) GetData() <-chan amqp.Delivery {

	c.logger.Log.Infoln("consumer has started")

	msgs, err := c.resultQueue.Consume(string(enum.ResultQueue))
	if err != nil {
		c.logger.Log.Fatalf("unable to open the consumer : %v", err)
	}

	return msgs
}

func (c *consumer) Operate(ctx context.Context, msgs <-chan amqp.Delivery) {

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				c.logger.Log.Infoln("consumer has closed")
				return
			}
			
			c.logger.Log.Infoln("consumer got data")

			var data queue.TaskStatus

			_ = json.Unmarshal(msg.Body, &data)

			
			err := c.db.TaskCompletedUpdateForAuth(ctx, data.Id)
			if err != nil {
				c.logger.Log.Errorf("error : unable to update the task status : %v", err)
			}

		case <-ctx.Done():
			return			
		}
	}

}