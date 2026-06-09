package workers

import (
	"context"
	"fmt"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/server"
	"orchestration_service/internal/types"
	"sync"
	"time"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type worker struct {
	id         int
	Jobs       <-chan types.Data
	PlanQueue  mq.QueueIface
	EmailQueue mq.QueueIface
	Db         *repository.Db
	logger     *logger.MyLogger
}

func StartWorkers(server server.Server) {

	for i := range server.NumberOfWorkers {
		worker := &worker{
			id:         i + 1,
			PlanQueue:  server.PlanQueue,
			EmailQueue: server.EmailQueue,
			Db:         server.Db,
			Jobs:       server.JobsChan,
			logger:     server.Logger,
		}

		server.WorkersWg.Add(1)

		go worker.start(server.WorkersWg)
	}

	server.Logger.Log.Infow("workers have started", zap.Int("number of workers", server.NumberOfWorkers))
}

func (w *worker) start(wg *sync.WaitGroup) {
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

			w.updateValueInDb(context.Background(), data)

		default:

			if data.NumberOfTries > 3 {
				w.numberOfTriesExceeded(context.Background(), data)
				continue
			}

			err := w.pushToQueue(context.Background(), data)
			if err != nil {
				w.exponentialBackOff(context.TODO(), data)
				continue
			}

			w.logger.Log.Infow("worker completed the task",
				zap.Int("worker id", w.id),
				zap.String("task name", data.Task),
			)
		}
	}
}

func (w *worker) pushToQueue(ctx context.Context, data types.Data) error {

	dataForSending := mqTypes.Data{
		DbId:          data.DbId,
		TaskName:      data.Task,
		Payload:       data.Payload,
		SentBy:        data.CreatedBy,
		TargetService: data.TargetService,
	}

	dataInBytes, err := dataForSending.ConvertIntoBytes()
	if err != nil {
		return fmt.Errorf("failed to convert data into bytes : %w", err)
	}

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

func (w *worker) exponentialBackOff(ctx context.Context, data types.Data) {

	numberOfTries := 3
	timeOut := time.Millisecond * 300

	for range numberOfTries {
		err := w.pushToQueue(ctx, data)
		if err != nil {
			time.Sleep(timeOut)
			timeOut += timeOut + 300
			continue
		}

		w.logger.Log.Infow("worker completed the task",
			zap.Int("worker id", w.id),
			zap.String("task name", data.Task),
		)

		return
	}

	switch data.CreatedBy {
	case enum.ServiceName_AUTH_SERVICE.String():
		err := w.Db.Auth.UpdateTaskStatusWithNumberOfTries(context.Background(), data.DbId, enum.TaskStatus_TASK_NOT_COMPLETED.String())
		if err != nil {
			w.logger.Log.Errorw("failed to update the task to failed in exponential backoff",
				zap.Int("worker id", w.id),
				zap.String("in service", data.CreatedBy),
				zap.Error(err),
			)
			return
		}
	case enum.ServiceName_TRACKER_SERVICE.String():
		err := w.Db.Tracker.UpdateTaskStatusWithNumberOfTries(context.Background(), data.DbId, enum.TaskStatus_TASK_NOT_COMPLETED.String())
		if err != nil {
			w.logger.Log.Errorw("failed to update the task to failed in exponential backoff",
				zap.Int("worker id", w.id),
				zap.String("in service", data.CreatedBy),
				zap.Error(err),
			)
			return
		}
	}

	w.logger.Log.Infow("failed to push the task to queue in exponential backoff and inserted data into push_to_queue_failed table",
		zap.Int("worker id", w.id),
		zap.String("task name", data.Task),
	)
}

func (w *worker) numberOfTriesExceeded(ctx context.Context, data types.Data) {
	switch data.CreatedBy {
	case enum.ServiceName_AUTH_SERVICE.String():
		err := w.Db.Auth.UpdateTaskStatus(ctx, data.DbId, enum.TaskStatus_TASK_FAILED.String())
		if err != nil {

			w.logger.Log.Errorw("failed to update the task to failed",
				zap.String("in service", data.CreatedBy),
				zap.Error(err),
			)
			return

		}

	case enum.ServiceName_TRACKER_SERVICE.String():
		err := w.Db.Tracker.UpdateTaskStatus(ctx, data.DbId, enum.TaskStatus_TASK_FAILED.String())
		if err != nil {

			w.logger.Log.Errorw("failed to update the task to failed",
				zap.String("in service", data.CreatedBy),
				zap.Error(err),
			)

			return
		}
	}

	w.logger.Log.Infow("worker updated the task to failed",
		zap.Int("worker id", w.id),
		zap.String("reason", "number of tries exceeded"),
		zap.String("task name", data.Task),
	)
}

func (w *worker) updateValueInDb(ctx context.Context, data types.Data) {
	switch data.TargetService {
	case enum.ServiceName_AUTH_SERVICE.String():

		err := w.Db.Auth.UpdateTaskStatus(ctx, data.DbId, data.Status)
		if err != nil {
			w.logger.Log.Errorw("worker falied to complete the task",
				zap.Int("worker_id", w.id),
				zap.String("task_name", data.Task),
				zap.String("created_by", data.CreatedBy),
				zap.String("targetServie", data.TargetService),
				zap.String("db index value", data.DbId),
				zap.Error(err),
			)
			return
		}

	case enum.ServiceName_TRACKER_SERVICE.String():

		err := w.Db.Tracker.UpdateTaskStatus(ctx, data.DbId, data.Status)
		if err != nil {
			w.logger.Log.Errorw("worker falied to complete the task",
				zap.Int("worker_id", w.id),
				zap.String("task_name", data.Task),
				zap.String("created_by", data.CreatedBy),
				zap.String("targetServie", data.TargetService),
				zap.String("db index value", data.DbId),
				zap.Error(err),
			)
			return
		}
	}

	w.logger.Log.Infow("worker completed the task",
		zap.Int("worker id", w.id),
		zap.String("task name", data.Task),
	)
}
