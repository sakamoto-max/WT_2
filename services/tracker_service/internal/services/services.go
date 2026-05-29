package services

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"tracker_service/internal/domain"

	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	trackerpb "github.com/sakamoto-max/wt_2_proto/shared/tracker"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	ErrWorkoutOngoing = errors.New("workout is ongoing")
)

func (s *Service) StartEmptyWorkout(ctx context.Context, in *trackerpb.StartEmptyWorkoutReq) (*trackerpb.GeneralResp, error) {
	trackerId, err := s.cache.TrackerId.GetTrackerId(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	if trackerId != "" {
		return nil, myerrors.BadReqErrMaker(ErrWorkoutOngoing)
	}

	r, err := s.planClient.GetEmptyPlanId(ctx, &planpb.SendUserID{UserId: in.UserId})
	if err != nil {
		return nil, fmt.Errorf("error getting data from plan server : %w", err)
	}

	if trackerId, err = s.pg.Start.StartWorkout(ctx, domain.StartWorkout{UserId: in.UserId, PlanId: r.EmptyPlanId}); err != nil {
		return nil, err
	}

	if err = s.cache.TrackerId.SetTrackerId(ctx, in.UserId, trackerId); err != nil {
		return nil, err
	}

	return &trackerpb.GeneralResp{
		Message: "an empty workout has started",
	}, nil
}

func (s *Service) StartWorkoutWithPlan(ctx context.Context, in *trackerpb.StartWorkoutWithPlanReq) (*trackerpb.StartWorkoutWithPlanResp, error) {

	trackerId, err := s.cache.TrackerId.GetTrackerId(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	if trackerId != "" {
		return nil, myerrors.BadReqErrMaker(ErrWorkoutOngoing)
	}

	r, err := s.planClient.GetPlanByName(ctx, &planpb.GetPlanByNameReq{UserId: in.UserId, PlanName: in.PlanName})
	if err != nil {
		return nil, fmt.Errorf("error getting data from plan server : %w", err)
	}

	err = s.cache.CurrentPlan.SetUserCurrentPlanName(ctx, in.UserId, in.PlanName)
	if err != nil {
		return nil, err
	}

	err = s.cache.Plan.SetPlanWithExercises(ctx, in.UserId, in.PlanName, &r.ExerciseNames)
	if err != nil {
		return nil, err
	}

	trackerId, err = s.pg.Start.StartWorkout(ctx, domain.StartWorkout{UserId: in.UserId, PlanId: r.PlanId})
	if err != nil {
		return nil, err
	}

	err = s.cache.TrackerId.SetTrackerId(ctx, in.UserId, trackerId)
	if err != nil {
		return nil, err
	}

	return &trackerpb.StartWorkoutWithPlanResp{
		Message:         fmt.Sprintf("workout with plan %v has started", in.PlanName),
		PlanName:        in.PlanName,
		ExercisesInPlan: r.ExerciseNames,
	}, nil
}

func (s *Service) EndWorkout(ctx context.Context, in *trackerpb.EndWorkoutReq) (*trackerpb.EndWorkoutResp, error) {

	trackerId, err := s.cache.TrackerId.GetTrackerId(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	if trackerId == "" {
		return nil, myerrors.BadReqErrMaker(fmt.Errorf("user doesn't have any workout ongoing"))
	}

	planName, err := s.cache.CurrentPlan.GetUserCurrentPlanName(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	var withOutbox bool

	var newExercisesPerformed *[]string

	data := domain.ConvertToLocal(in)

	for i := range len(data.Workout) {

		exerciseName := data.Workout[i].ExerciseName

		in := exerpb.SendExerciseName{
			UserId:       in.UserId,
			ExerciseName: exerciseName,
		}

		resp, err := s.exerClient.ExerciseExistsReturnId(ctx, &in)
		if err != nil {

			return nil, err
		}

		a := resp.ExerciseId

		data.Workout[i].ExerciseId = a
	}

	var TriggerDbCommit bool

	switch {
	case planName == "":
		TriggerDbCommit = true
	case planName != "":
		allExersInPlan, err := s.cache.Plan.GetPlanWithExercises(ctx, in.UserId, planName)
		if err != nil {
			return nil, err
		}

		conflictLevel, err := s.cache.Conflict.GetConflictLevel(ctx, in.UserId)
		if err != nil {
			return nil, err
		}

		switch conflictLevel {
		case 0:

			exercisesInTracker := data.GetAllExercises()

			err := s.cache.TrackerData.SetUserTrackerData(ctx, in.UserId, data)
			if err != nil {
				return nil, err
			}

			resp := checkExercisesNotPerformed(allExersInPlan, exercisesInTracker)
			if resp != nil {
				err := s.cache.Conflict.SetConflictLevel(ctx, in.UserId, 1)
				if err != nil {
					return nil, err
				}

				return &trackerpb.EndWorkoutResp{
					RequestStatus:   resp.RequestStatus,
					Message:         resp.Message,
					Reason:          resp.Reason.Error(),
					ExerciseNames:   resp.ExerciseNames,
					ConflictOccured: true,
				}, nil
			}

			resp = checkIfNewExercisesAdded(allExersInPlan, exercisesInTracker)
			if resp != nil {
				err := s.cache.Conflict.SetConflictLevel(ctx, in.UserId, 2)
				if err != nil {
					return nil, err
				}

				err = s.cache.NewExercises.SetUserNewExercises(ctx, in.UserId, &resp.ExerciseNames)
				if err != nil {
					return nil, err
				}

				newExercisesPerformed = &resp.ExerciseNames
				//
				return &trackerpb.EndWorkoutResp{
					RequestStatus:   resp.RequestStatus,
					Message:         resp.Message,
					Reason:          resp.Reason.Error(),
					ExerciseNames:   resp.ExerciseNames,
					ConflictOccured: true,
				}, nil
			}

			TriggerDbCommit = true

		case 1:

			yes := data.UserResponse

			if !yes {
				return &trackerpb.EndWorkoutResp{
					Message: "please continue the workout",
				}, nil
			}

			data, err = s.cache.TrackerData.GetUserTrackerData(ctx, in.UserId)
			if err != nil {
				return nil, err
			}

			exercisesInTracker := data.GetAllExercises()

			resp := checkIfNewExercisesAdded(allExersInPlan, exercisesInTracker)
			if resp != nil {
				err := s.cache.Conflict.SetConflictLevel(ctx, in.UserId, 2)
				if err != nil {
					return nil, err
				}

				err = s.cache.NewExercises.SetUserNewExercises(ctx, in.UserId, &resp.ExerciseNames)
				if err != nil {
					return nil, err
				}

				newExercisesPerformed = &resp.ExerciseNames

				return nil, resp
			}

			TriggerDbCommit = true

		case 2:
			yes := data.UserResponse

			data, err = s.cache.TrackerData.GetUserTrackerData(ctx, in.UserId)
			if err != nil {
				return nil, err
			}

			if yes {
				exerciseNames, err := s.cache.NewExercises.GetUserNewExercises(ctx, in.UserId)
				if err != nil {
					return nil, err
				}

				newExercisesPerformed = exerciseNames

				withOutbox = true
			}

			TriggerDbCommit = true
		}
	}

	if TriggerDbCommit {

		switch withOutbox {
		case true:
			err := s.pg.End.EndWorkoutWithOutbox(ctx, in.UserId, trackerId, data, planName, newExercisesPerformed)
			if err != nil {
				return nil, err
			}
		case false:
			err := s.pg.End.EndWorkout(ctx, trackerId, data)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := s.cache.UserData.DelAllUserData(ctx, in.UserId, planName); err != nil {
		return nil, err
	}

	return &trackerpb.EndWorkoutResp{
		Message: "workout ended successfully",
	}, nil
}

func (s *Service) PING(ctx context.Context, in *trackerpb.PingTrackReq) (*trackerpb.PingTrackResp, error) {
	r := trackerpb.PingTrackResp{}

	return &r, nil
}

func (s *Service) GetHealth(ctx context.Context, in *trackerpb.GetHealthReq) (*trackerpb.GetHealthResp, error) {
	pgRespTime := s.pg.Metrics.GetRespTime(ctx)
	redisRespTime := s.cache.Metrics.GetRespTime(ctx)

	if pgRespTime == nil && redisRespTime == nil {
		return &trackerpb.GetHealthResp{
			PostgresRespTime: nil,
			RedisRespTime:    nil,
		}, nil
	}
	if redisRespTime == nil {
		return &trackerpb.GetHealthResp{
			PostgresRespTime: durationpb.New(*pgRespTime),
			RedisRespTime:    nil,
		}, nil
	}
	if pgRespTime == nil {
		return &trackerpb.GetHealthResp{
			PostgresRespTime: nil,
			RedisRespTime:    durationpb.New(*redisRespTime),
		}, nil
	}

	return &trackerpb.GetHealthResp{
		PostgresRespTime: durationpb.New(*pgRespTime),
		RedisRespTime:    durationpb.New(*redisRespTime),
	}, nil
}

func (s *Service) CancelWorkout(ctx context.Context, in *trackerpb.CancelWorkoutReq) (*trackerpb.CancelWorkoutResp, error) {
	trackerId, err := s.cache.TrackerId.GetTrackerId(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	if trackerId == "" {
		return nil, myerrors.BadReqErrMaker(fmt.Errorf("user has no workout ongoing"))
	}

	planName, err := s.cache.CurrentPlan.GetUserCurrentPlanName(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	switch {
	case planName == "":
		err := s.cache.TrackerId.DelTrackerId(ctx, in.UserId)
		if err != nil {
			return nil, err
		}

	case planName != "":
		err := s.cache.UserData.DelAllUserData(ctx, in.UserId, planName)
		if err != nil {
			return nil, err
		}
	}

	err = s.pg.Cancel.DeleteTrackerIdInPG(ctx, trackerId)
	if err != nil {
		return nil, err
	}
	return &trackerpb.CancelWorkoutResp{
		Message: "workout has been successfully canceled",
	}, nil
}

func checkIfNewExercisesAdded(originalPlanExers *[]string, userSent *[]string) *myerrors.Conflict {

	var newExercises []string

	for _, eachExer := range *userSent {
		exists := slices.Contains(*originalPlanExers, eachExer)

		if !exists {
			newExercises = append(newExercises, eachExer)
		}
	}

	if len(newExercises) != 0 {
		resp := myerrors.Conflict{
			RequestStatus: "INCOMPLETE",
			Reason:        myerrors.ErrNewExercises,
			Message:       "update plan?",
			ExerciseNames: newExercises,
		}

		return &resp
	}

	return nil
}

func checkExercisesNotPerformed(originalPlanExers *[]string, userSent *[]string) *myerrors.Conflict {
	var notPerformed []string

	for _, eachExer := range *originalPlanExers {
		exists := slices.Contains(*userSent, eachExer)

		if !exists {
			notPerformed = append(notPerformed, eachExer)
		}
	}

	if len(notPerformed) > 0 {

		resp := myerrors.Conflict{
			RequestStatus: "INCOMPLETE",
			Reason:        myerrors.ErrNotPerformed,
			Message:       "still complete?",
			ExerciseNames: notPerformed,
		}

		return &resp
	}

	return nil
}
