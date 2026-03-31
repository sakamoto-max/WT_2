package services

import (
	"context"
	"fmt"
	"time"
	myerrors "wt/pkg/my_errors"
	"wt/pkg/enum"

	// customerrors "tracker_service/internal/custom_errors"
	// "tracker_service/internal/models"
	"tracker_service/internal/models"
	"tracker_service/internal/repository"

	// "tracker_service/internal/user"
	exerpb "workout-tracker/proto/shared/exercise"
	planpb "workout-tracker/proto/shared/plan"
	// "github.com/redis/go-redis/v9"
)

type Service struct {
	Db      *repository.DBs
	PClient planpb.PlanServiceClient
	EClient exerpb.ExerciseServiceClient
}

func (s *Service) GetHealth(ctx context.Context) (*time.Duration, *time.Duration) {

	// check resp time of pg

	pgRespTime := s.Db.GetPostgresRespTime(ctx)
	redisRespTime := s.Db.GetRedisRespTime(ctx)

	return pgRespTime, redisRespTime
}

func NewService(Db *repository.DBs, planClient planpb.PlanServiceClient, exerClient exerpb.ExerciseServiceClient) *Service {
	return &Service{Db: Db, PClient: planClient, EClient: exerClient}
}

func (s *Service) StartEmptyWorkoutSer(ctx context.Context, userID string) error {
	// get empty plan_id of user

	ongoing, err := s.Db.CheckIfWorkoutIsOngoing(ctx, userID)
	if err != nil {
		return err
	}

	if ongoing {
		return myerrors.ErrWorkoutOngoing
	}

	r, err := s.PClient.GetPlanByName(ctx, &planpb.GetPlanByNameReq{UserId: userID, PlanName: string(enum.EmptyPlanName)})
	if err != nil {
		return fmt.Errorf("error getting data from plan server : %w", err)
	}
	
	trackerId, err := s.Db.StartWorkout(ctx, userID, r.PlanId)
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

func (s *Service) StartWorkoutWithPlanSer(ctx context.Context, userId string, planName string) (*[]string, error) {
	// check if plan Name exists
	// if exists get the plan_id

	ongoing, err := s.Db.CheckIfWorkoutIsOngoing(ctx, userId)
	if err != nil {
		return nil, err
	}

	if ongoing {
		return nil, myerrors.ErrWorkoutOngoing
	}

	r, err := s.PClient.GetPlanByName(ctx, &planpb.GetPlanByNameReq{UserId: userId, PlanName: planName})
	if err != nil {
		return nil, fmt.Errorf("error getting data from plan server : %w", err)
	}

	trackerId, err := s.Db.StartWorkout(ctx, userId, r.PlanId)
	if err != nil {
		return nil, err
	}

	err = s.Db.SetTrackerId(ctx, userId, trackerId)
	if err != nil {
		err := s.Db.RevertStartWorkout(ctx, trackerId)
		if err != nil {
			return nil, err
		}
	}

	return &r.ExerciseNames, nil
}

func (s *Service) EndWorkoutSer(ctx context.Context, userId string, data *models.Tracker) error {

	// get tracker ID from redis
	trackerId, err := s.Db.GetTrackerId(ctx, userId)
	if err != nil {
		return err
	}
	// do the db ops
	err = s.Db.EndWorkout(ctx, trackerId, data)
	if err != nil {
		return err
	}
	// del the tracker ID
	err = s.Db.DelTrackerId(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}
