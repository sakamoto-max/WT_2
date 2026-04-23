package workers

import (
	"context"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/types"
	"sync"
	// "wt/pkg/enum"
	"github.com/sakamoto-max/wt_2-pkg/enum"
	// "wt/pkg/logger"
	"github.com/sakamoto-max/wt_2-pkg/logger"
	// mq "wt/pkg/queue"
	mq "github.com/sakamoto-max/wt_2-pkg/queue"

	"go.uber.org/zap"
	// "golang.org/x/text/cases"
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
			return
		}

		dataInBytes, _ := data.ConvertToBytes()

		err := w.PushToQueue(ctx, dataInBytes, data.TargetService, data.Task)
		if err != nil {

			err := w.Db.TaskNotCompleted(ctx, data.TargetService, data.Id)
			if err != nil {
				w.logger.Log.Errorw(
					"error in updating the task to not completed",
					zap.Int("worker_id", w.id),
					zap.Error(err),
				)
				return
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
	case string(enum.PlanService):
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

	case string(enum.EmailService):
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