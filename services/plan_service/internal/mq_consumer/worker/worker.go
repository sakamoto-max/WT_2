package worker

import (
	"context"
	"fmt"
	"plan_service/internal/mq_consumer/types"
	"plan_service/internal/repository"
	"sync"
	mq "github.com/sakamoto-max/rabbit_mq/queue"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"go.uber.org/zap"
)

type worker struct {
	id          int
	db          *repository.DBs
	logger      *logger.MyLogger
	jobs        <-chan types.Data
	resultQueue *mq.MessageQueue
	exerclient  exerpb.ExerciseServiceClient
}

func newWorker(id int, repo *repository.DBs, logger *logger.MyLogger, jobs <-chan types.Data, resQueue *mq.MessageQueue, client exerpb.ExerciseServiceClient) *worker {
	return &worker{
		id:          id,
		db:          repo,
		logger:      logger,
		jobs:        jobs,
		resultQueue: resQueue,
		exerclient:  client,
	}
}

func MakeWorkers(numberOfWorkers int, repo *repository.DBs, logger *logger.MyLogger, jobs <-chan types.Data, resQueue *mq.MessageQueue, Client exerpb.ExerciseServiceClient) []*worker {

	var workers []*worker

	for i := 1; i <= numberOfWorkers; i++ {
		w := newWorker(i, repo, logger, jobs, resQueue, Client)

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
		case enum.TaskName_CREATE_EMPTY_PLAN_FOR_USER.String():

			userId, _ := msg.GetUserId()

			err := w.db.CreateEmptyPlan(context.TODO(), userId)
			if err != nil {
				w.logger.Log.Infow("failed to create empty plan for the user", zap.Int("worker", w.id), zap.Error(err))

				w.SendDataToResQ(
					msg.DbId,
					enum.ServiceName_PLAN_SERVICE.String(),
					enum.ServiceName_AUTH_SERVICE.String(),
					msg.Task,
					enum.TaskStatus_TASK_NOT_COMPLETED.String(),
				)

				continue
			}

			err = w.SendDataToResQ(
				msg.DbId,
				enum.ServiceName_PLAN_SERVICE.String(),
				enum.ServiceName_AUTH_SERVICE.String(),
				msg.Task,
				enum.TaskStatus_TASK_COMPLETED.String(),
			)
			if err != nil {

			}

		case enum.TaskName_UPDATE_PLAN.String():

			userId, err := msg.GetUserId()
			if err != nil {
				w.logger.Log.Errorln(err)
			}
			planName, err := msg.GetPlanName()
			if err != nil {
				w.logger.Log.Errorln(err)
			}
			newExercises, err := msg.GetNewExercises()
			if err != nil {
				w.logger.Log.Errorln(err)
			}

			planId, err := w.db.ReturnsPlanId(context.TODO(), userId, planName)
			if err != nil {
				w.logger.Log.Infow("failed to get plan_id for the user", zap.Int("worker", w.id), zap.Error(err))

				err := w.SendDataToResQ(
					msg.DbId,
					enum.ServiceName_PLAN_SERVICE.String(),
					enum.ServiceName_TRACKER_SERVICE.String(),
					msg.Task,
					enum.TaskStatus_TASK_NOT_COMPLETED.String(),
				)
				if err != nil {
					w.logger.Log.Errorw("error sending data to the result queue", zap.Error(err))
				}
				continue
			}

			var exerciseIds []string

			for _, exerciseName := range newExercises {
				in := exerpb.SendExerciseName{
					UserId:       userId,
					ExerciseName: exerciseName,
				}
				resp, err := w.exerclient.ExerciseExistsReturnId(context.TODO(), &in)
				if err != nil {
					w.logger.Log.Errorw("error occured while getting the exercise id", zap.Int("worker_id", w.id), zap.Error(err))
				}

				exerciseIds = append(exerciseIds, resp.ExerciseId)
			}

			// w.exerclient.ExerciseExistsReturnId(context.TODO(), )

			err = w.db.AddExercisesToPlan(context.TODO(), planId, &exerciseIds)
			if err != nil {
				w.logger.Log.Infow("failed to update plan for the user", zap.Int("worker", w.id), zap.Error(err))
				w.SendDataToResQ(
					msg.DbId,
					enum.ServiceName_PLAN_SERVICE.String(),
					enum.ServiceName_TRACKER_SERVICE.String(),
					msg.Task,
					enum.TaskStatus_TASK_NOT_COMPLETED.String(),
				)

				continue
			}

			err = w.SendDataToResQ(
				msg.DbId,
				enum.ServiceName_PLAN_SERVICE.String(),
				enum.ServiceName_TRACKER_SERVICE.String(),
				msg.Task,
				enum.TaskStatus_TASK_COMPLETED.String(),
			)
			if err != nil {
				w.logger.Log.Errorw("error occured while sending data to the res Queue", zap.Error(err))
			}
		}
	}
}

func (w *worker) SendDataToResQ(id string, sentBy string, targetService string, taskName string, taskStatus string) error {

	fmt.Println("msg.id", id)

	d := mq.NewTaskStatus(
		id,
		sentBy,
		targetService,
		taskName,
		taskStatus,
	)

	dataInBytes := d.ConvertToBytes()

	err := w.resultQueue.Publish(context.TODO(), dataInBytes)
	if err != nil {
		return fmt.Errorf("failed to publish to resQueue : %w", err)
	}

	return nil
}
