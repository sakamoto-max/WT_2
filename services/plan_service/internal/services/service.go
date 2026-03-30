package services

import (
	"context"
	"fmt"
	// customerrors "plan_service/internal/custom_errors"
	"plan_service/internal/models"
	"plan_service/internal/repository"
	"strings"
	"time"
	exerpb "workout-tracker/proto/shared/exercise"
	// myerrors "wt/pkg/my_errors"
)

type Service struct {
	Db      *repository.DBs
	GClient exerpb.ExerciseServiceClient
}

func NewService(Db *repository.DBs, grpcCli exerpb.ExerciseServiceClient) *Service {
	return &Service{Db: Db, GClient: grpcCli}
}

// DONE
func (s *Service) CreatePlan(ctx context.Context, userId int, planName string, exerciseNames *[]string) error {
	var exerciseIDs []string
	PlanName := strings.ToLower(planName)

	for _, exerciseName := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: exerciseName})
		if err != nil {
			return err
		}

		exerciseIDs = append(exerciseIDs, r.ExerciseId)
	}

	err := s.Db.CreatePlan(ctx, userId, PlanName, exerciseIDs)
	if err != nil {
		return err
	}

	return nil
}

// DONE
func (s *Service) GetAllPlansSer(ctx context.Context, userId int) (int, *[]models.Plan2, error) {

	var allPlans []models.Plan2

	planNamesWithIds, err := s.Db.GetAllUserPlans(ctx, userId)
	if err != nil {
		return 0, nil, err
	}

	for _, eachPlan := range *planNamesWithIds {

		var plan models.Plan2
		var exerciseNames []string

		exeriseIDs, err := s.Db.GetAllExercisesByPlanID(ctx, eachPlan.Id)
		if err != nil {
			return 0, nil, err
		}

		for _, id := range *exeriseIDs {
			exerciseName, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: id, UserId: int64(userId)})
			if err != nil {
				return 0, nil, err
			}

			exerciseNames = append(exerciseNames, exerciseName.ExerciseName)
		}

		plan.PlanName = eachPlan.PlanName
		plan.Exercises = exerciseNames
		allPlans = append(allPlans, plan)
	}

	return len(*planNamesWithIds), &allPlans, nil

}

// DONE
func (s *Service) GetPlanByNameSer(ctx context.Context, userId int, planName string) (string, *[]string, error) {

	planId, err := s.Db.ReturnsPlanId(ctx, userId, planName)
	if err != nil {
		return "", nil, err
	}

	exerciseIds, err := s.Db.GetAllExercisesByPlanID(ctx, planId)
	if err != nil {
		return "", nil, err
	}

	var allExercises []string
	for _, exerciseId := range *exerciseIds {

		r, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: exerciseId, UserId: int64(userId)})
		if err != nil {
			return "", nil, err
		}

		allExercises = append(allExercises, r.ExerciseName)
	}

	return planName, &allExercises, nil
}

// DONE
func (s *Service) AddExercisesToPlan(ctx context.Context, userId int, planName string, exerciseNames *[]string) (*models.Plan2, error) {

	
	planId, err := s.Db.ReturnsPlanId(ctx, userId, planName)
	if err != nil {
		return nil, err
	}

	var exerciseIds []string
	for _, eachName := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: eachName, UserId: int64(userId)})
		if err != nil{
			return nil, err
		}

		exerciseIds = append(exerciseIds, r.ExerciseId)
	}

	err = s.Db.AddExercisesToPlan(ctx, planId, &exerciseIds)
	if err != nil {
		return nil, err
	}

	resp, err := MakeRespForAddingNewExer(ctx, userId, planId, planName, s)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DONE
func (s *Service) DeleteExerciseFromPlan(ctx context.Context, userId int, planName string, exerciseNames *[]string) (*models.Plan2, error) {

	planId, err := s.Db.ReturnsPlanId(ctx, userId, planName)
	if err != nil {
		return nil, err
	}
	
	var exerciseIds []string

	for _, v := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: v, UserId: int64(userId)})
		if err != nil {
			return nil, fmt.Errorf("error getting data from exercise server : %w", err)
		}

		exerciseIds = append(exerciseIds, r.ExerciseId)
	}

	err = s.Db.DeleteExerciseFromPlan(ctx, planId, &exerciseIds)
	if err != nil {
		return nil, err
	}

	resp, err := MakeRespForAddingNewExer(ctx, userId, planId, planName, s)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DONE
func (s *Service) DeletePlanSer(ctx context.Context, userId int, planName string) error {

	planId, err := s.Db.ReturnsPlanId(ctx, userId, planName)
	if err != nil {
		return err
	}

	err = s.Db.DeletePlan(ctx, userId, planId)
	if err != nil {
		return err
	}

	return nil
}


// func (s *Service) PlanExistsReturnId(ctx context.Context, userId int, planName string) (bool, int, error) {

// 	planId, err := s.Db.ReturnsPlanId(ctx, userId, planName)
// 	if err != nil{
// 		return err
// 	}


// 	// return s.Db.PlanExistsReturnsId(ctx, userID, planName)
// }
// func (s *Service) PlanExistsReturnPlan(ctx context.Context, userId int, planName string) (bool, int, *[]int, error) {
// 	return s.Db.PlanExistsReturnPlan(ctx, userId, planName)
// }

// func (s *Service) GetEmptyPlanId(ctx context.Context, userId int) (int, error) {
// 	return s.Db.GetEmptyPlanID(ctx, userId)
// }
// func (s *Service) CreateEmptyPlan(ctx context.Context, userId int) error {
// 	return s.Db.CreateEmptyPlan(ctx, userId)
// }

func MakeRespForAddingNewExer(ctx context.Context, userId int, planId string, planName string, s *Service) (*models.Plan2, error) {

	var resp models.Plan2
	var allExercises []string

	exerciseIds, err := s.Db.GetAllExercisesByPlanID(ctx, planId)
	if err != nil {
		return &resp, err
	}

	for _, exerciseId := range *exerciseIds {

		r, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: exerciseId, UserId: int64(userId)})
		if err != nil {
			return &resp, fmt.Errorf("error getting data from exercise grpc server : %w", err)
		}

		allExercises = append(allExercises, r.ExerciseName)
	}

	resp.PlanName = planName
	resp.Exercises = allExercises

	return &resp, nil
}

func (s *Service) GetHealth(ctx context.Context) (*time.Duration, *time.Duration) {

	// check resp time of pg

	pgRespTime := s.Db.GetPostgresRespTime(ctx)
	redisRespTime := s.Db.GetRedisRespTime(ctx)

	return pgRespTime, redisRespTime
}
