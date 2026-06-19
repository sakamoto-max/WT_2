package workers

import (
	"context"
	"fmt"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/server"
	"orchestration_service/internal/types"
	"sync"

	mq "github.com/sakamoto-max/rabbit_mq/queue"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type worker struct {
	id              int
	Jobs            chan types.Data
	PlanQueue       mq.QueueIface
	EmailQueue      mq.QueueIface
	DeadLetterQueue mq.QueueIface
	Db              *repository.Db
	logger          *logger.MyLogger
}

func StartWorkersForFetcher(server server.Server) {

	for i := range server.NumberOfWorkers {
		worker := &worker{
			id:              i + 1,
			PlanQueue:       server.PlanQueue,
			EmailQueue:      server.EmailQueue,
			DeadLetterQueue: server.DeadLetterQueue,
			Db:              server.Db,
			Jobs:            server.FetcherJobsChan,
			logger:          server.Logger,
		}

		server.FetcherWorkersWg.Add(1)

		go worker.start(server.FetcherWorkersWg)
	}

	server.Logger.Log.Infow("fetcher workers have started", zap.Int("number of workers", server.NumberOfWorkers))
}
func StartWorkersForConsumer(server server.Server) {

	for i := range server.NumberOfWorkers {
		worker := &worker{
			id:              i + 1,
			PlanQueue:       server.PlanQueue,
			EmailQueue:      server.EmailQueue,
			DeadLetterQueue: server.DeadLetterQueue,
			Db:              server.Db,
			Jobs:            server.ConsumerJobsChan,
			logger:          server.Logger,
		}

		server.ConsumerWorkerWg.Add(1)

		go worker.start(server.ConsumerWorkerWg)
	}

	server.Logger.Log.Infow("consumer workers have started", zap.Int("number of workers", server.NumberOfWorkers))
}

func (w *worker) start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		data, ok := <-w.Jobs
		if !ok {
			return
		}

		skip, err := w.Db.Cache.SkipTask(context.Background(), data)
		if err != nil {
			w.logger.Log.Errorw("failed to check if the task should be skipped or not",
				zap.Int("worker id", w.id),
				zap.Error(err),
			)
			continue
		}

		if skip {
			w.Jobs <- data
			continue
		}

		w.logger.Log.Infow("worker received a task",
			zap.Int("worker id", w.id),
			zap.String("task name", data.Task),
		)

		if data.WorkersTries > 5 {
			err := w.numberOfTriesExceeded(context.Background(), data)
			if err != nil {
				w.logger.Log.Errorw("failed to update the task to failed",
					zap.String("in service", data.CreatedBy),
					zap.Error(err),
				)
				continue
			}

			w.logger.Log.Infow("worker updated the task to failed",
				zap.Int("worker id", w.id),
				zap.String("reason", "number of tries exceeded"),
				zap.String("task name", data.Task),
			)
			continue
		}

		switch data.Task {
		case enum.TaskName_UPDATE_VALUE_IN_DB.String():

			err := w.updateValueInDb(context.Background(), data)
			if err != nil {
				w.logger.Log.Errorw("worker falied to complete the task",
					zap.Int("worker_id", w.id),
					zap.String("task_name", data.Task),
					zap.String("created_by", data.CreatedBy),
					zap.String("targetServie", data.TargetService),
					zap.String("db index value", data.DbId),
					zap.Error(err),
				)
				continue
			}

			w.logger.Log.Infow("worker completed the task",
				zap.Int("worker id", w.id),
				zap.String("task name", data.Task),
			)

			continue

		default:

			err := w.pushToQueue(context.Background(), data)
			if err != nil {
				err := w.exponentialBackOff(context.TODO(), data)
				if err != nil {
					w.logger.Log.Errorw("failed to set task timeout", zap.Int("worker_id", w.id), zap.Error(err))
					continue
				}

				w.logger.Log.Infow("worker set the task timeout", zap.Int("worker_id", w.id))
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

func (w *worker) exponentialBackOff(ctx context.Context, data types.Data) error {

	data.WorkersTries = data.WorkersTries + 1

	switch data.TimeOut.Level {
	case 0:
		data.InitTimeOut()
	default:
		data.IncreaseTimeOut()
	}

	err := w.Db.Cache.SetTaskTimeOut(ctx, data)
	if err != nil {
		return err
	}

	w.Jobs <- data

	return nil
}

func (w *worker) numberOfTriesExceeded(ctx context.Context, data types.Data) error {

	dataInBytes, _ := data.ConvertToBytes()

	err := w.DeadLetterQueue.Publish(ctx, dataInBytes)
	if err != nil {
		return err
	}

	return nil
}

func (w *worker) updateValueInDb(ctx context.Context, data types.Data) error {
	switch data.TargetService {
	case enum.ServiceName_AUTH_SERVICE.String():

		err := w.Db.Auth.UpdateTaskStatus(ctx, data.DbId, data.Status)
		if err != nil {
			return err

		}

	case enum.ServiceName_TRACKER_SERVICE.String():

		err := w.Db.Tracker.UpdateTaskStatus(ctx, data.DbId, data.Status)
		if err != nil {
			return err
		}
	}

	return nil
}

