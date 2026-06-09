package consumer

import (
	"email_service/internals/server"
	"email_service/internals/types"
	"encoding/json"

	"github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

type consumer struct {
	emailQueue *queue.MessageQueue
	logger     *logger.MyLogger
	jobs       chan<- types.Data
}

func StartConsumer(server server.Server) {

	consumer := &consumer{
		emailQueue: server.EmailQueue,
		logger:     server.Logger,
		jobs:       server.JobsChan,
	}

	go consumer.StartListening()

	server.Logger.Log.Infoln("consumer has started")

}

func (c *consumer) StartListening() {

	msgs, err := c.emailQueue.Consume()
	if err != nil {
		c.logger.Log.Fatalw("failed to consume from email queue", zap.Error(err))
	}

	for {
		msg, ok := <-msgs

		if !ok {
			return
		}

		data := convertIntoStruct(&msg.Body)

		c.logger.Log.Infow("consumer has got data", 
			zap.String("targer service", data.TargetService), 	
			zap.String("task name", data.TaskName),
		)

		c.jobs <- data
	}

}

func convertIntoStruct(data *[]byte) types.Data {

	var D types.Data

	_ = json.Unmarshal(*data, &D)

	return D
}
