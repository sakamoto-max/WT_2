package producer

import (
	"context"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"

	"go.uber.org/zap"
)

type Producer struct {
	logger *logger.MyLogger
	resQueue *queue.MessageQueue
}

func NewProducer(logger *logger.MyLogger, resQueue *queue.MessageQueue) *Producer {
	return &Producer{
		logger: logger,
		resQueue: resQueue,
	}
}

func (p *Producer) TaskFailed(dBId string, originatedBy string, taskName string) {

	targerService := enum.ServiceName_EMAIL_SERVICE.String()
	dbUpdateValue := enum.TaskStatus_TASK_NOT_COMPLETED.String()

	d := queue.NewTaskStatus(dBId, targerService, originatedBy, taskName, dbUpdateValue)
	dataInBytes := d.ConvertToBytes()

	err := p.resQueue.Publish(context.TODO(), dataInBytes)
	if err != nil{
		p.logger.Log.Errorw(
			"failed to push to result queue",
			zap.Error(err),
		)
		return
	}

	p.logger.Log.Infow("pushed data to the result queue")
}
func (p *Producer) TaskCompleted(dBId string, originatedBy string, taskName string) {

	targerService := enum.ServiceName_EMAIL_SERVICE.String()
	dbUpdateValue := enum.TaskStatus_TASK_COMPLETED.String()

	d := queue.NewTaskStatus(dBId, targerService, originatedBy, taskName, dbUpdateValue)
	dataInBytes := d.ConvertToBytes()

	err := p.resQueue.Publish(context.TODO(), dataInBytes)
	if err != nil{
		p.logger.Log.Errorw(
			"failed to push to result queue",
			zap.Error(err),
		)
	}

	p.logger.Log.Infow("pushed data to the result queue")
}