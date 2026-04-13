package services

import (
	"context"
	"fmt"
	"time"
	"tracker_service/internal/models"
	"tracker_service/internal/repository"
	exerpb "workout-tracker/proto/shared/exercise"
	planpb "workout-tracker/proto/shared/plan"
	"wt/pkg/enum"
	myerrors "wt/pkg/my_errors"
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

	err = s.Db.SetUserWorkingOutWithPlan(ctx, userID, false)
	if err != nil {
		return err
	}

	trackerId, err := s.Db.StartWorkout(ctx, userID, r.PlanId)
	if err != nil {
		return err
	}

	err = s.Db.SetTrackerIdAndOngoingWorkout(ctx, userID, trackerId)
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

	err = s.Db.SetPlanWithExercises(ctx, userId, planName, &r.ExerciseNames)
	if err != nil {
		return nil, err
	}

	// set user_id:%v:current_workout_plan_name

	err = s.Db.SetUserWorkingOutWithPlan(ctx, userId, true)
	if err != nil {
		return nil, err
	}

	trackerId, err := s.Db.StartWorkout(ctx, userId, r.PlanId)
	if err != nil {
		return nil, err
	}

	err = s.Db.SetTrackerIdAndOngoingWorkout(ctx, userId, trackerId)
	if err != nil {
		err := s.Db.RevertStartWorkout(ctx, trackerId)
		if err != nil {
			return nil, err
		}
	}

	return &r.ExerciseNames, nil
}

func (s *Service) EndWorkoutSer(ctx context.Context, userId string, data *models.Tracker) error {

	yes, err := s.Db.GetUserWorkingOutWithPlan(ctx, userId)
	if err != nil {
		return err
	}

	if yes {
		// s.Db.GetPlanWithExercises(ctx, userId)
		// check if all the exercises in the workout are performed
		// what is performed
		// atleast one set should be performed
	}

	for i := range len(data.Workout) {

		exerciseName := data.Workout[i].ExerciseName

		in := exerpb.SendExerciseName{
			UserId:       userId,
			ExerciseName: exerciseName,
		}

		resp, err := s.EClient.ExerciseExistsReturnId(ctx, &in)
		if err != nil {
			return err
		}

		a := resp.ExerciseId

		data.Workout[i].ExerciseId = a
	}

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
	err = s.Db.DelTrackerIdAndOngoingWorkout(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func exercisesNotPerformed(data *models.Tracker, exerciseNames *[]string) {

}

