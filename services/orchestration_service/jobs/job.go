package jobs

import (
	"context"
	"orchestration_service/repository"
	"orchestration_service/types"
	"wt/pkg/enum"
	mq "wt/pkg/queue"
)

func OperateForPlan(ctx context.Context, DbId string, v types.Data, planQueue *mq.MessageQueue, Db *repository.DB) error {
	dataInBytes, err := v.ConvertToBytes()
	if err != nil {
		return err
	}

	err = planQueue.Publish(ctx, dataInBytes, string(enum.ApplicationJsonType))
	if err != nil {
		err := Db.TaskNotCompletedForAuth(ctx, DbId)
		if err != nil{
			return err
		}
		return err
	}

	return nil
}

func OperateForEmail(ctx context.Context, v types.Data, emailQueue *mq.MessageQueue, Db *repository.DB) error {
	dataInBytes, _ := v.ConvertToBytes()

	err := emailQueue.Publish(ctx, dataInBytes, string(enum.ApplicationJsonType))
	if err != nil {
		return err
	}
	err = Db.TaskPendingUpdateForAuth(ctx, v.Id)
	if err != nil {
		return err
	}

	return nil

}
