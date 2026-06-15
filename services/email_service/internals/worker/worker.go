package worker

import (
	"email_service/internals/server"
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

func StartWorkers(server server.Server) {

	for i := range server.NumberOfWorkers {
		worker := &worker{
			id:         i + 1,
			logger:     server.Logger,
			service:    server.Service,
			jobs:       server.JobsChan,
			senderChan: server.SenderChan,
		}

		server.WorkersWg.Add(1)
		go worker.Start(server.WorkersWg)
	}

	server.Logger.Log.Infow("workers have started", zap.Int("number of workers", server.NumberOfWorkers))
}

func (w *worker) Start(wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		msg, ok := <-w.jobs

		if !ok {
			return
		}

		w.logger.Log.Infow("worker received task", 
			zap.Int("worker_id", w.id), 
			zap.String("task", msg.TaskName),
		)

		switch msg.TaskName {
		case enum.TaskName_SEND_EMAIL_FOR_SIGNING_UP.String():

			email, err := msg.GetEmail()
			if err != nil {
				w.senderChan <- msg.Failed(err)
				continue
			}

			err = w.service.SendWelcomeEmail(email)
			if err != nil {
				w.senderChan <- msg.Failed(err)
				continue
			}

			w.senderChan <- msg.Succeded()
		}

	}
}
