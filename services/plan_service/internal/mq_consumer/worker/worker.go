package worker

import (
	"context"
	"plan_service/internal/repository"
	"sync"
	"wt/pkg/enum"
	"wt/pkg/logger"
	"wt/pkg/queue"
	mq "wt/pkg/queue"
	"wt/pkg/types"

	"go.uber.org/zap"
)

// job of the worker
// 1. read from jobs
// 2. do the tasks

type worker struct {
	id          int
	db          *repository.DBs
	logger      *logger.MyLogger
	jobs        <-chan types.Data
	resultQueue *mq.MessageQueue
}

func newWorker(id int, repo *repository.DBs, logger *logger.MyLogger, jobs <-chan types.Data, resQueue *mq.MessageQueue) *worker {
	return &worker{
		id:          id,
		db:          repo,
		logger:      logger,
		jobs:        jobs,
		resultQueue: resQueue,
	}
}

func MakeWorkers(numberOfWorkers int, repo *repository.DBs, logger *logger.MyLogger, jobs <-chan types.Data, resQueue *mq.MessageQueue) []*worker {

	var workers []*worker

	for i := 1; i <= numberOfWorkers; i++ {
		w := newWorker(i, repo, logger, jobs, resQueue)

		workers = append(workers, w)
	}

	return workers
}

func (w *worker) Work(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		msg, ok := <-w.jobs

		if !ok {
			w.logger.Log.Infow("worker stopped", zap.Int("worker_id", w.id))
			return
		}

		w.logger.Log.Infow(
			"worker received a job",
			zap.Int("worker_id", w.id),
			zap.String("task", msg.Task),
		)

		switch msg.Task {
		case string(enum.CreateEmptyPlanForUser):
			err := w.db.CreateEmptyPlan(context.TODO(), msg.Payload["user_id"])
			if err != nil {
				w.logger.Log.Infow("failed to create empty plan for the user", zap.Int("worker", w.id), zap.Error(err))
				w.UpdateTaskFailed(msg.Id, string(enum.AuthService), string(enum.PlanService), string(enum.CreateEmptyPlanForUser))
				continue
			}

			w.UpdateTaskCompleted(
				msg.Id,
				string(enum.PlanService),
				string(enum.AuthService),
				string(enum.CreateEmptyPlanForUser),
			)

		}
	}
}

func (w *worker) UpdateTaskCompleted(id string, targerService string, orginatedBy string, taskName string) {

	d := queue.NewTaskStatus(
		id,
		targerService,
		orginatedBy,
		taskName,
		string(enum.TaskCompleted),
	)

	dataInBytes := d.ConvertToBytes()

	err := w.resultQueue.Publish(context.TODO(), dataInBytes, string(enum.ApplicationJsonType))
	if err != nil {
		w.logger.Log.Infof("failed to publish to reqQueue", zap.Int("worker", w.id), zap.Error(err))
	}
}

func (w *worker) UpdateTaskFailed(id string, targerService string, orginatedBy string, taskName string) {

	d := queue.NewTaskStatus(
		id,
		targerService,
		orginatedBy,
		taskName,
		string(enum.TaskNotCompleted),
	)

	dataInBytes := d.ConvertToBytes()

	err := w.resultQueue.Publish(context.TODO(), dataInBytes, string(enum.ApplicationJsonType))
	if err != nil {
		w.logger.Log.Infof("failed to publish to reqQueue", zap.Int("worker", w.id), zap.Error(err))
	}
}
