package services

import (
	"context"
	"fmt"
	myerrors "wt/pkg/my_errors"

	// customerrors "tracker_service/internal/custom_errors"
	// "tracker_service/internal/models"
	"tracker_service/internal/repository"
	"tracker_service/internal/user"
	exerpb "workout-tracker/proto/shared/exercise"
	planpb "workout-tracker/proto/shared/plan"
)

type Service struct {
	Db      *repository.DBs
	PClient planpb.PlanServiceClient
	EClient exerpb.ExerciseServiceClient
}

func NewService(Db *repository.DBs, planClient planpb.PlanServiceClient, exerClient exerpb.ExerciseServiceClient) *Service {
	return &Service{Db: Db, PClient: planClient, EClient: exerClient}
}

func (s *Service) StartEmptyWorkoutSer(ctx context.Context, userID int) error {
	// get empty plan_id of user

	ongoing, err := s.Db.CheckIfWorkoutIsOngoing(ctx, userID)
	if err != nil{
		return err
	}

	if ongoing {
		return myerrors.ErrWorkoutOngoing
	}

	

	r, err := s.PClient.GetEmptyPlanId(ctx, &planpb.SendUserID{UserId: int64(userID)})
	if err != nil {
		return fmt.Errorf("error getting data from plan server : %w", err)
	}

	trackerId, err := s.Db.StartWorkout(ctx, userID, int(r.EmptyPlanId))
	if err != nil {
		return err
	}

	err = s.Db.SetTrackerId(ctx, userID, trackerId)
	if err != nil {
		err := s.Db.RevertStartWorkout(ctx, trackerId)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (s *Service) StartWorkoutWithPlanSer(ctx context.Context, userId int, planName string) (*[]string, error) {
	// check if plan Name exists
	// if exists get the plan_id

	var allExerNames []string
	// var resp models.Plan

	r, err := s.PClient.PlanExistsReturnPlan(ctx, &planpb.SendPlanName{UserId: int64(userId), PlanName: planName})
	if err != nil{
		return &allExerNames, fmt.Errorf("error getting data from plan server : %w", err)
	}

	if !r.Exists {
		return &allExerNames, fmt.Errorf("plan doesnt exist")
	}
	
	for _, v := range r.ExerciseIds {
		r, err := s.EClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: v})
		if err != nil{
			return &allExerNames, fmt.Errorf("error getting data from exercise server : %w", err)
		}

		allExerNames = append(allExerNames, r.ExerciseName)
	}

	// do db operations
	trackerId, err := s.Db.StartWorkout(ctx, userId, int(r.PlanId))
	if err != nil{
		return &allExerNames, err
	}

	err = s.Db.SetTrackerId(ctx, userId, trackerId)
	if err != nil{
		err := s.Db.RevertStartWorkout(ctx, trackerId)
		if err != nil{
			return &allExerNames, err
		}
	}

	return &allExerNames, nil
}

func (s *Service) EndWorkoutSer(ctx context.Context, userId int, data *user.Tracker) error {

	// get tracker ID from redis
	trackerId, err := s.Db.GetTrackerId(ctx, userId)
	if err != nil{
		return err
	}
	// do the db ops
	err = s.Db.EndWorkout(ctx, trackerId, data)
	if err != nil{
		return err
	}
	// del the tracker ID
	err = s.Db.DelTrackerId(ctx, userId)
	if err != nil{
		return err
	}
	
	return nil
}
