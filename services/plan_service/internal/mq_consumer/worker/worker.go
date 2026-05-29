package worker

import (
	"context"
	"plan_service/internal/client"
	"plan_service/internal/domain"
	"plan_service/internal/mq_consumer/types"
	"plan_service/internal/repository"
	"sync"

	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"go.uber.org/zap"
)

type worker struct {
	id         int
	db         *repository.Db
	logger     *logger.MyLogger
	jobsChan   <-chan types.Data
	senderChan chan<- mqTypes.Data
	exerclient client.ExerClientIface
}

func MakeWorkers(numberOfWorkers int, repo *repository.Db, logger *logger.MyLogger, jobs <-chan types.Data, senderChan chan<- mqTypes.Data, Client client.ExerClientIface) []*worker {

	var workers []*worker

	for i := range numberOfWorkers {

		w := &worker{
			id:         i + 1,
			db:         repo,
			logger:     logger,
			jobsChan:   jobs,
			senderChan: senderChan,
			exerclient: Client,
		}

		workers = append(workers, w)
	}

	return workers
}

func (w *worker) Work(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		msg, ok := <-w.jobsChan

		if !ok {
			return
		}

		w.logger.Log.Infow("worker received a job", zap.Int("worker_id", w.id), zap.String("task", msg.TaskName))

		switch msg.TaskName {
		case enum.TaskName_CREATE_EMPTY_PLAN_FOR_USER.String():

			userId, err := msg.GetUserId()
			if err != nil {
				w.logger.Log.Errorln(err)
				continue
			}

			err = w.db.PlanCommandRepo.CreateEmptyPlan(context.TODO(), userId)
			if err != nil {

				w.senderChan <- mqTypes.Data{
					DbId:          msg.DbId,
					TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
					SentBy:        msg.TargetService,
					TaskStatus:    enum.TaskStatus_TASK_FAILED.String(),
					TargetService: msg.SentBy,
					Err:           err,
				}
				continue
			}

			w.senderChan <- mqTypes.Data{
				DbId:          msg.DbId,
				TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
				SentBy:        msg.TargetService,
				TaskStatus:    enum.TaskStatus_TASK_COMPLETED.String(),
				TargetService: msg.SentBy,
			}

		case enum.TaskName_UPDATE_PLAN.String():

			userId, err := msg.GetUserId()
			if err != nil {
				w.logger.Log.Errorln(err)
				continue
			}

			planName, err := msg.GetPlanName()
			if err != nil {
				w.logger.Log.Errorln(err)
				continue
			}

			newExercises, err := msg.GetNewExercises()
			if err != nil {
				w.logger.Log.Errorln(err)
				continue
			}

			planId, err := w.db.PlanQueryRepo.GetPlanId(context.TODO(), domain.GetPlan{UserId: userId, PlanName: planName})
			if err != nil {
				w.senderChan <- mqTypes.Data{
					DbId:          msg.DbId,
					TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
					SentBy:        msg.TargetService,
					TaskStatus:    enum.TaskStatus_TASK_FAILED.String(),
					TargetService: msg.SentBy,
					Err:           err,
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

					w.senderChan <- mqTypes.Data{
						DbId:          msg.DbId,
						TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
						SentBy:        msg.TargetService,
						TaskStatus:    enum.TaskStatus_TASK_FAILED.String(),
						TargetService: msg.SentBy,
						Err:           err,
					}
				}

				exerciseIds = append(exerciseIds, resp.ExerciseId)
			}

			err = w.db.PlanExericseRepo.AddExercisesToPlan(context.TODO(), planId, &exerciseIds)
			if err != nil {

				w.senderChan <- mqTypes.Data{
					DbId:          msg.DbId,
					TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
					SentBy:        msg.TargetService,
					TaskStatus:    enum.TaskStatus_TASK_FAILED.String(),
					TargetService: msg.SentBy,
					Err:           err,
				}
				continue

			}

			w.senderChan <- mqTypes.Data{
				DbId:          msg.DbId,
				TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
				SentBy:        msg.TargetService,
				TaskStatus:    enum.TaskStatus_TASK_COMPLETED.String(),
				TargetService: msg.SentBy,
			}
		}
	}
}
