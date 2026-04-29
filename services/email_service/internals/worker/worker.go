package worker

import (
	"email_service/internals/producer"
	"email_service/internals/services"
	"email_service/internals/types"
	"sync"

	// "github.com/sakamoto-max/wt_2_pkg/enum"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/sakamoto-max/wt_2_pkg/logger"

	// "github.com/sakamoto-max/wt_2_pkg/types"

	"go.uber.org/zap"
)

type worker struct {
	id      int
	logger  *logger.MyLogger
	service *services.Service
	jobs <-chan types.Data
	producer *producer.Producer
}

func MakeWorkers(NumberOfWorkers int, logger *logger.MyLogger, service *services.Service, jobs <-chan types.Data, producer *producer.Producer) []*worker {

	var workers []*worker
	for i := 1; i <= NumberOfWorkers; i++ {
		w := newWorker(i, logger, service, jobs, producer)
		workers = append(workers, w)
	}

	return workers

}
func newWorker(id int, logger *logger.MyLogger, service *services.Service, jobs <-chan types.Data, producer *producer.Producer) *worker {
	return &worker{
		id:      id,
		logger:  logger,
		service: service,
		jobs: jobs,
		producer: producer,
	}
}


func (w *worker) Work(wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		msg, ok := <- w.jobs
	
		if !ok {
			w.logger.Log.Infow("worker is stopping", zap.Int("worker_id", w.id))
			return
		}


		w.logger.Log.Infow("worker received task", zap.Int("worker_id", w.id), zap.String("task", msg.Task))

		switch msg.Task{
		case enum.TaskName_SEND_EMAIL_FOR_SIGNING_UP.String():

			email, err := msg.GetEmail()
			if err != nil {
				w.logger.Log.Errorw(
					"error getting email",
					zap.Int("worker_id", w.id),
					zap.Error(err),
				)
			}

			err = w.service.SendWelcomeEmail(email)
			if err != nil{
				w.producer.TaskFailed(msg.DbId, enum.ServiceName_AUTH_SERVICE.String(), msg.Task)
				w.logger.Log.Errorw("error occured while sending email", zap.Int("worker_id", w.id), zap.String("email", email), zap.Error(err))
			}
		}

		w.producer.TaskCompleted(msg.DbId, enum.ServiceName_AUTH_SERVICE.String(), msg.Task)
		w.logger.Log.Infow(
			"task completed", 
			zap.Int("worker_id", w.id), 
			zap.String("task", msg.Task),
		)
	}



}