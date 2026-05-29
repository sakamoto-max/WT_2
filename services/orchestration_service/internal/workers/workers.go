package workers

import (
	"context"
	"fmt"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/types"
	"sync"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type worker struct {
	id         int
	Jobs       <-chan types.Data
	PlanQueue  *mq.MessageQueue
	EmailQueue *mq.MessageQueue
	Db         *repository.Db
	logger     *logger.MyLogger
}

func MakeWorkers(NumberOfWorkers int, planQueue *mq.MessageQueue, emailQueue *mq.MessageQueue, db *repository.Db, jobs <-chan types.Data, logger *logger.MyLogger) []*worker {

	var workers []*worker

	for i := 1; i <= NumberOfWorkers; i++ {
		w := &worker{
			id:         i,
			PlanQueue:  planQueue,
			EmailQueue: emailQueue,
			Db:         db,
			Jobs:       jobs,
			logger:     logger,
		}
		workers = append(workers, w)
	}
	return workers
}

func (w *worker) Work(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		data, ok := <-w.Jobs
		if !ok {
			return
		}

		w.logger.Log.Infow("worker received a task",
			zap.Int("worker id", w.id),
			zap.String("task name", data.Task),
		)

		switch data.Task {
		case enum.TaskName_UPDATE_VALUE_IN_DB.String():
			switch data.TargetService {
			case enum.ServiceName_AUTH_SERVICE.String():

				fmt.Println("status", data.Status)

				err := w.Db.Auth.UpdateTaskStatus(ctx, data.DbId, data.Status)
				if err != nil {
					w.logger.Log.Errorw("failed to update task status", zap.Int("worker_id", w.id), zap.String("task_name", data.Task), zap.String("created_by", data.CreatedBy), zap.String("targetServie", data.TargetService), zap.String("db index value", data.DbId), zap.Error(err))
				}

			case enum.ServiceName_TRACKER_SERVICE.String():

				err := w.Db.Tracker.UpdateTaskStatus(ctx, data.DbId, data.Status)
				if err != nil {
					w.logger.Log.Errorw("failed to update task status", zap.Int("worker_id", w.id), zap.String("task_name", data.Task), zap.String("created_by", data.CreatedBy), zap.String("targetServie", data.TargetService), zap.String("db index value", data.DbId), zap.Error(err))
				}
			}

			w.logger.Log.Infow("worker completed the task",
				zap.Int("worker id", w.id),
				zap.String("task name", data.Task),
			)

		default:

			if data.NumberOfTries > 3 {

				switch data.CreatedBy {
				case enum.ServiceName_AUTH_SERVICE.String():
					err := w.Db.Auth.UpdateTaskStatus(ctx, data.DbId, enum.TaskStatus_TASK_FAILED.String())
					if err != nil {

						w.logger.Log.Errorw("failed to update the task to failed",
							zap.String("in service", data.CreatedBy),
							zap.Error(err),
						)

					}

				case enum.ServiceName_TRACKER_SERVICE.String():
					err := w.Db.Tracker.UpdateTaskStatus(ctx, data.DbId, enum.TaskStatus_TASK_FAILED.String())
					if err != nil {

						w.logger.Log.Errorw("failed to update the task to failed",
							zap.String("in service", data.CreatedBy),
							zap.Error(err),
						)
					}
				}
				continue
			}

			err := w.PushToQueue(ctx, data)
			if err != nil {


				switch data.CreatedBy {
				case enum.ServiceName_AUTH_SERVICE.String():
					err := w.Db.Auth.UpdateTaskStatusWithNumberOfTries(ctx, data.DbId, enum.TaskStatus_TASK_PENDING.String())
					if err != nil {
						w.logger.Log.Errorw("failed to update the task to failed", zap.String("in service", data.CreatedBy), zap.Error(err))
					}

				case enum.ServiceName_TRACKER_SERVICE.String():
					err := w.Db.Tracker.UpdateTaskStatusWithNumberOfTries(ctx, data.DbId, enum.TaskStatus_TASK_PENDING.String())
					if err != nil {
						w.logger.Log.Errorw("failed to update the task to failed", zap.String("in service", data.CreatedBy), zap.Error(err))
					}
				}

				continue
			}
		}
	}
}

func (w *worker) PushToQueue(ctx context.Context, data types.Data) error {

	dataForSending := mqTypes.Data{
		DbId:          data.DbId,
		TaskName:      data.Task,
		Payload:       data.Payload,
		SentBy:        data.CreatedBy,
		TargetService: data.TargetService,
	}

	dataInBytes, _ := dataForSending.ConvertIntoBytes()

	switch data.TargetService {
	case enum.ServiceName_PLAN_SERVICE.String():

		err := w.PlanQueue.Publish(ctx, dataInBytes)
		if err != nil {
			return fmt.Errorf("failed to publish data to the plan queue : %w", err)
		}

	case enum.ServiceName_EMAIL_SERVICE.String():
		err := w.EmailQueue.Publish(ctx, dataInBytes)
		if err != nil {
			return fmt.Errorf("failed to publish data to the email queue : %w", err)
		}
	}
	return nil
}
