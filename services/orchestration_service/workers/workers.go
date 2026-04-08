package workers

import (
	"context"
	"orchestration_service/jobs"
	"orchestration_service/repository"
	"orchestration_service/types"
	"sync"
	"wt/pkg/enum"
	"wt/pkg/logger"
	mq "wt/pkg/queue"

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

func (w *worker) DoWork(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		data, ok := <-w.Jobs

		if !ok {
			return
		}

		switch data.TargetService {

		case string(enum.PlanService):
			err := jobs.OperateForPlan(ctx, data.Id, data, w.PlanQueue, w.Db)
			if err != nil {
				w.logger.Log.Errorw(
					"error in completing the operation",
					zap.Int("worker_id", w.id),
					zap.String("task", data.Task),
					zap.String("targer_service", data.TargetService),
					zap.Error(err),
				)
				continue
			}

			w.logger.Log.Infow(
				"worker completed the task",
				zap.Int("worker_id", w.id),
				zap.String("task", data.Task),
				zap.String("target_service", data.TargetService),
			)

		case string(enum.EmailService):
			err := jobs.OperateForEmail(ctx, data, w.EmailQueue, w.Db)
			if err != nil {
				w.logger.Log.Errorw(
					"error in completing the operation",
					zap.Int("worker_id", w.id),
					zap.String("task", data.Task),
					zap.String("targer_service", data.TargetService),
					zap.Error(err),
				)

				continue
			}

			w.logger.Log.Infow(
				"worker completed the task",
				zap.Int("worker_id", w.id),
				zap.String("task", data.Task),
				zap.String("target_service", data.TargetService),
			)

		}
	}
}
