package services

import (
	"context"
	"tracker_service/internal/repository"
)

type Service struct {
	Db *repository.DBs
}

func NewService(Db *repository.DBs) *Service {
	return &Service{Db: Db}
}


func (s *Service) StartEmptyWorkoutSer(ctx context.Context, userID int) (error) {
	// get empty plan_id of user
	planID := 1

	trackerId, err := s.Db.StartWorkout(ctx, userID, planID)
	if err != nil{
		return err

	}
	err = s.Db.SetTrackerId(ctx, userID, trackerId)
	if err != nil{
		err := s.Db.RevertStartWorkout(ctx, trackerId)
		if err != nil{
			return err
		}
		return err
	}

	return nil
}

func StartWorkoutWithPlanSer(ctx context.Context, userId int, planName string) {
	// check if plan Name exists
	// if exists get the plan_id

	// if not 

	
	


}

func EndWorkoutSer() {
	
}

