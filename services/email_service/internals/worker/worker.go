package worker

import (
	"email_service/internals/services"
	"email_service/internals/types"
	"sync"

	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"go.uber.org/zap"
)

type worker struct {
	id         int
	logger     *logger.MyLogger
	service    *services.Service
	jobs       <-chan types.Data
	senderChan chan<- types.Data
}

func MakeWorkers(NumberOfWorkers int, logger *logger.MyLogger, service *services.Service, jobs <-chan types.Data, senderChan chan<- types.Data) []*worker {

	var workers []*worker

	for i := range NumberOfWorkers {
		w := &worker{
			id:         i + 1,
			logger:     logger,
			service:    service,
			jobs:       jobs,
			senderChan: senderChan,
		}

		workers = append(workers, w)

	}
	return workers
}

func (w *worker) Work(wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		msg, ok := <-w.jobs

		if !ok {
			return
		}

		w.logger.Log.Infow("worker received task", zap.Int("worker_id", w.id), zap.String("task", msg.TaskName))

		switch msg.TaskName {
		case enum.TaskName_SEND_EMAIL_FOR_SIGNING_UP.String():

			email, err := msg.GetEmail()
			if err != nil {
				w.senderChan <- types.Data{
					DbId:          msg.DbId,
					TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
					Status:        enum.TaskStatus_TASK_FAILED.String(),
					SentBy:        msg.TargetService,
					TargetService: msg.SentBy,
					Err:           err,
				}
				continue
			}

			err = w.service.SendWelcomeEmail(email)
			if err != nil {
				w.senderChan <- types.Data{
					DbId:          msg.DbId,
					TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
					Status:        enum.TaskStatus_TASK_FAILED.String(),
					SentBy:        msg.TargetService,
					TargetService: msg.SentBy,
					Err:           err,
				}
				continue
			}

			w.senderChan <- types.Data{
				DbId:          msg.DbId,
				TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
				Status:        enum.TaskStatus_TASK_COMPLETED.String(),
				SentBy:        msg.TargetService,
				TargetService: msg.SentBy,
				Err:           err,
			}
		}

	}
}

// consumer responsibilities :
// 1. listens to the email queue
// 2. if it gets any data -> send it to jobs queue

// consumer depends on -> email queue, jobs chan

// worker responsibilities :
// 1. listen to the jobs chan
// 2. if data is received -> perform the TaskName
// 3. send the result to the producer queue

// worker depends on -> jobs chan, producer chan, email service

// producer
// 1. lisents to the producer queue
// 2. if data is received it will push it to the res queue

// producer depends on -> producer chan, resQueue
