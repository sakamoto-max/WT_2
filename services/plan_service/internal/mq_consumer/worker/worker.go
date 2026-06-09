package worker

import (
	"context"
	"fmt"
	"plan_service/internal/client"
	"plan_service/internal/domain"
	"plan_service/internal/mq_consumer/server"
	"plan_service/internal/mq_consumer/types"
	"plan_service/internal/repository"
	"plan_service/internal/repository/cache"
	"sync"

	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"go.uber.org/zap"
)

type worker struct {
	id         int
	cache      *cache.Cache
	db         *repository.Db
	logger     *logger.MyLogger
	jobsChan   <-chan types.Data
	senderChan chan<- mqTypes.Data
	exerclient client.ExerClientIface
}

func StartWorkers(server server.Server) {

	for i := range server.NumberOfWorkers {

		worker := &worker{
			id:         i + 1,
			db:         server.Db,
			logger:     server.Logger,
			jobsChan:   server.JobsChan,
			senderChan: server.SendersChan,
			exerclient: server.ExerClient,
			cache: server.Cache,
		}

		server.WorkerWg.Add(1)
		go worker.Work(server.WorkerWg)
	}

	server.Logger.Log.Infow("workers have started", zap.Int("number of workers", server.NumberOfWorkers))

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

			fmt.Println("worker got all the data")
			
			planId, err := w.db.PlanQueryRepo.GetPlanId(context.TODO(), domain.GetPlan{UserId: userId, PlanName: planName})
			if err != nil {
				fmt.Print("worker got an error %w", err)
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
			fmt.Println("worker got the plan id")
			
			var exerciseIds []string
			
			for _, exerciseName := range newExercises {
				in := exerpb.SendExerciseName{
					UserId:       userId,
					ExerciseName: exerciseName,
				}
				resp, err := w.exerclient.ExerciseExistsReturnId(context.TODO(), &in)
				if err != nil {
					fmt.Print("worker got an error %w", err)
					
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
			
			fmt.Println("worker got all the exercise ids")
			
			err = w.db.PlanExericseRepo.AddExercisesToPlan(context.TODO(), planId, &exerciseIds)
			if err != nil {
				fmt.Print("worker got an error %w", err)
				
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
			fmt.Println("worker completed adding exercises to the plan")
			
			err = w.cache.UserPlan.DelUserPlan(context.TODO(), domain.GetPlan{UserId: userId, PlanName: planName})
			if err != nil {
				fmt.Print("worker got an error %w", err)
				w.senderChan <- mqTypes.Data{
					DbId:          msg.DbId,
					TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
					SentBy:        msg.TargetService,
					TaskStatus:    enum.TaskStatus_TASK_FAILED.String(),
					TargetService: msg.SentBy,
					Err:           err,
				}
			}
			
			fmt.Println("worker deleted the user plan in cache ")


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
