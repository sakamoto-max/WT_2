package workers

import (
	"context"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/types"
	"sync"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	mq "github.com/sakamoto-max/rabbit_mq/queue" 
	"go.uber.org/zap"
)

type worker struct {
	id         int
	PlanQueue  *mq.MessageQueue
	EmailQueue *mq.MessageQueue
	Db         *repository.DB
	Jobs       <-chan types.Data
	logger     *logger.MyLogger
}

func MakeWorkers(NumberOfWorkers int, planQueue *mq.MessageQueue, emailQueue *mq.MessageQueue, db *repository.DB, jobs <-chan types.Data, logger *logger.MyLogger) []*worker {

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
			w.logger.Log.Infow("worker stopped", zap.Int("id", w.id), zap.String("reason", "shutdown"))
			return
		}

		w.logger.Log.Infow("worker received a task", zap.Int("worker_id", w.id), zap.String("task_name", data.Task), zap.String("created_by", data.CreatedBy))

		if data.NumberOfTries > 3 {

			err := w.Db.TaskFailed(ctx, data.CreatedBy, data.DbId)
			if err != nil {
				w.logger.Log.Errorw(
					"failed to update task to failed",
					zap.Error(err),
				)
			}
		}

		dataInBytes, _ := data.ConvertToBytes()

		err := w.PushToQueue(ctx, dataInBytes, data.TargetService, data.Task)
		if err != nil {
			err := w.Db.TaskNotCompletedUpdateTries(ctx, data.CreatedBy, data.DbId)
			if err != nil {
				w.logger.Log.Errorw(
					"error in updating the task to not completed",
					zap.Int("worker_id", w.id),
					zap.Error(err),
				)
				continue
			}
		}

		w.logger.Log.Infow(
			"worker successfully pushed to the queue",
			zap.Int("worker_id", w.id),
			zap.String("task", data.Task),
			zap.String("target_service", data.TargetService),
		)
	}
}

func (w *worker) PushToQueue(ctx context.Context, data *[]byte, targetService string, task string) error {

	switch targetService {
	case enum.ServiceName_PLAN_SERVICE.String():
		err := w.PlanQueue.Publish(ctx, data)
		if err != nil {
			w.logger.Log.Errorw(
				"error in completing the operation",
				zap.Int("worker_id", w.id),
				zap.String("task", task),
				zap.String("targer_service", targetService),
				zap.Error(err),
			)
			return err
		}
	case enum.ServiceName_EMAIL_SERVICE.String():
		err := w.EmailQueue.Publish(ctx, data)
		if err != nil {
			w.logger.Log.Errorw(
				"error in completing the operation",
				zap.Int("worker_id", w.id),
				zap.String("task", task),
				zap.String("targer_service", targetService),
				zap.Error(err),
			)
			return err
		}
	}
	return nil
}
