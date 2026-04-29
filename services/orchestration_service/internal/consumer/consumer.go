package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"orchestration_service/internal/repository"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	mq "github.com/sakamoto-max/rabbit_mq/queue"
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
	
	msgs, err := c.resultQueue.Consume(enum.QueueName_RESULT_QUEUE.String())
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

			var data mq.TaskStatus

			_ = json.Unmarshal(msg.Body, &data)

			fmt.Printf("target_service in consumer is %v", data.TargetService)
			
			err := c.db.TaskCompletedUpdate(ctx, data.TargetService, data.Id)
			if err != nil {
				c.logger.Log.Errorf("error : unable to update the task status : %v", err)
			}
		case <-ctx.Done():
			c.logger.Log.Infoln("consumer is closing")
			return			
		}
	}
}