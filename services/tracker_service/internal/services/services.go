package services

import (
	"context"
	"fmt"
	"slices"
	"time"
	"tracker_service/internal/models"
	"tracker_service/internal/repository"
	exerpb "github.com/sakamoto-max/wt_2-proto/shared/exercise"
	planpb "github.com/sakamoto-max/wt_2-proto/shared/plan"
	myerrors "github.com/sakamoto-max/wt_2-pkg/my_errors"
)

type Service struct {
	db      *repository.DBs
	pClient planpb.PlanServiceClient
	eClient exerpb.ExerciseServiceClient
}

func (s *Service) GetHealth(ctx context.Context) (*time.Duration, *time.Duration) {

	pgRespTime := s.db.GetPostgresRespTime(ctx)
	redisRespTime := s.db.GetRedisRespTime(ctx)

	return pgRespTime, redisRespTime
}

func NewService(Db *repository.DBs, planClient planpb.PlanServiceClient, exerClient exerpb.ExerciseServiceClient) *Service {
	return &Service{db: Db, pClient: planClient, eClient: exerClient}
}

func (s *Service) StartEmptyWorkoutSer(ctx context.Context, userID string) error {
	// get empty plan_id of user

	trackerId, err := s.db.GetTrackerId(ctx, userID)
	if err != nil {
		return err
	}

	if trackerId != "" {
		return myerrors.ErrWorkoutOngoing
	}

	r, err := s.pClient.GetEmptyPlanId(ctx, &planpb.SendUserID{UserId: userID})
	// r, err := s.pClient.GetEmptyPlanId().GetPlanByName(ctx, &planpb.GetPlanByNameReq{UserId: userID, PlanName: string(enum.EmptyPlanName)})
	if err != nil {
		return fmt.Errorf("error getting data from plan server : %w", err)
	}


	trackerId, err = s.db.StartWorkout(ctx, userID, r.EmptyPlanId)
	if err != nil {
		return err
	}

	err = s.db.SetTrackerId(ctx, userID, trackerId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) StartWorkoutWithPlanSer(ctx context.Context, userId string, planName string) (*[]string, error) {

	trackerId, err := s.db.GetTrackerId(ctx, userId)
	if err != nil {
		return nil, err
	}

	if trackerId != "" {
		return nil, myerrors.ErrWorkoutOngoing
	}

	r, err := s.pClient.GetPlanByName(ctx, &planpb.GetPlanByNameReq{UserId: userId, PlanName: planName})
	if err != nil {
		return nil, fmt.Errorf("error getting data from plan server : %w", err)
	}

	err = s.db.SetUserCurrentPlanName(ctx, userId, planName)
	if err != nil {
		return nil, err
	}

	err = s.db.SetPlanWithExercises(ctx, userId, planName, &r.ExerciseNames)
	if err != nil {
		return nil, err
	}

	trackerId, err = s.db.StartWorkout(ctx, userId, r.PlanId)
	if err != nil {
		return nil, err
	}

	err = s.db.SetTrackerId(ctx, userId, trackerId)
	if err != nil {
		return nil, err
	}

	return &r.ExerciseNames, nil
}

func (s *Service) EndWorkoutSer(ctx context.Context, userId string, data *models.Tracker) (*string, error) {

	trackerId, err := s.db.GetTrackerId(ctx, userId)
	if err != nil {
		return nil, err
	}

	if trackerId == "" {
		return nil, myerrors.BadReqErrMaker(fmt.Errorf("user doesn't have any workout ongoing"))
	}

	planName, err := s.db.GetUserCurrentPlanName(ctx, userId)
	if err != nil {
		return nil, err
	}

	var withOutbox bool

	var newExercisesPerformed *[]string

	data, err = getExerIdOfEachExercise(ctx, data, userId, s)
	if err != nil {
		return nil, err
	}

	var TriggerDbCommit bool

	switch {
	case planName == "":
		TriggerDbCommit = true
	case planName != "":
		allExersInPlan, err := s.db.GetPlanWithExercises(ctx, userId, planName)
		if err != nil {
			return nil, err
		}

		conflictLevel, err := s.db.GetConflictLevel(ctx, userId)
		if err != nil {
			return nil, err
		}

		switch conflictLevel {
		case 0:

			exercisesInTracker := data.GetAllExercises()

			err := s.db.SetUserTrackerData(ctx, userId, data)
			if err != nil {
				return nil, err
			}

			resp := checkExercisesNotPerformed(allExersInPlan, exercisesInTracker)
			if resp != nil {
				err := s.db.SetConflictLevel(ctx, userId, 1)
				if err != nil {
					return nil, err
				}

				return nil, resp
			}

			resp = checkIfNewExercisesAdded(allExersInPlan, exercisesInTracker)
			if resp != nil {
				err := s.db.SetConflictLevel(ctx, userId, 2)
				if err != nil {
					return nil, err
				}

				err = s.db.SetUserNewExercises(ctx, userId, &resp.ExerciseNames)
				if err != nil {
					return nil, err
				}

				newExercisesPerformed = &resp.ExerciseNames

				return nil, resp
			}

			TriggerDbCommit = true

		case 1:

			yes := data.UserResponse

			if !yes {
				resp := "please continue the workout"
				return &resp, nil
			}

			data, err = s.db.GetUserTrackerData(ctx, userId)
			if err != nil {
				return nil, err
			}

			exercisesInTracker := data.GetAllExercises()

			resp := checkIfNewExercisesAdded(allExersInPlan, exercisesInTracker)
			if resp != nil {
				err := s.db.SetConflictLevel(ctx, userId, 2)
				if err != nil {
					return nil, err
				}

				// set to redis

				err = s.db.SetUserNewExercises(ctx, userId, &resp.ExerciseNames)
				if err != nil {
					return nil, err
				}

				newExercisesPerformed = &resp.ExerciseNames

				return nil, resp
			}

			TriggerDbCommit = true

		case 2:
			yes := data.UserResponse

			data, err = s.db.GetUserTrackerData(ctx, userId)
			if err != nil {
				return nil, err
			}

			if yes {
				// get newExercises list from redis
				exerciseNames, err := s.db.GetUserNewExercises(ctx, userId)
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

		trackerId, err := s.db.GetTrackerId(ctx, userId)
		if err != nil {
			return nil, err
		}

		switch withOutbox {
		case true:
			err := s.db.EndWorkoutWithOutbox(ctx, userId, trackerId, data, planName, newExercisesPerformed)
			if err != nil {
				return nil, err
			}
		case false:
			err := s.db.EndWorkout(ctx, trackerId, data)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := s.db.DelAllUserData(ctx, userId, planName); err != nil {
		return nil, err
	}

	return nil, nil
}

// helpers

func getExerIdOfEachExercise(ctx context.Context, data *models.Tracker, userId string, s *Service) (*models.Tracker, error) {
	for i := range len(data.Workout) {

		exerciseName := data.Workout[i].ExerciseName

		in := exerpb.SendExerciseName{
			UserId:       userId,
			ExerciseName: exerciseName,
		}

		resp, err := s.eClient.ExerciseExistsReturnId(ctx, &in)
		if err != nil {
			return nil, err
		}

		a := resp.ExerciseId

		data.Workout[i].ExerciseId = a
	}

	return data, nil
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

func (s *Service) CancelWorkout(ctx context.Context, userID string) error {

	trackerId, err := s.db.GetTrackerId(ctx, userID)
	if err != nil {
		return err
	}

	if trackerId == "" {
		return myerrors.BadReqErrMaker(fmt.Errorf("user has no workout ongoing"))
	}

	planName, err := s.db.GetUserCurrentPlanName(ctx, userID)
	if err != nil {
		return err
	}

	switch {
	case planName == "":
		err := s.db.DelTrackerId(ctx, userID)
		if err != nil {
			return err
		}

	case planName != "":
		err := s.db.DelAllUserData(ctx, userID, planName)
		if err != nil {
			return err
		}
	}

	err = s.db.DeleteTrackerIdInPG(ctx, trackerId)
	if err != nil {
		return err
	}
	return nil
}

// trackerId
// workoutOngoing
// workoutWithPlan
