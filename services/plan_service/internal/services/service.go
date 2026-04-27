package services

import (
	"context"
	"fmt"

	// customerrors "plan_service/internal/custom_errors"
	"plan_service/internal/models"
	"plan_service/internal/repository"
	"strings"
	"time"

	// exerpb "workout-tracker/proto/shared/exercise"
	exerpb "github.com/sakamoto-max/wt_2-proto/shared/exercise"
	// myerrors "wt/pkg/my_errors"
)

type ServiceIface interface {
	CreatePlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) error
	GetAllPlansSer(ctx context.Context, userId string) (int, *[]models.Plan2, error)
	GetPlanByNameSer(ctx context.Context, userId string, planName string) (string, string, *[]string, error)
	AddExercisesToPlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error)
	DeleteExerciseFromPlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error)
	DeletePlanSer(ctx context.Context, userId string, planName string) error
	GetEmptyPlanId(ctx context.Context, userId string) (string, error)
	GetHealth(ctx context.Context) (*time.Duration, *time.Duration)
}

type service struct {
	Db      *repository.DBs
	GClient exerpb.ExerciseServiceClient
}

func NewService(Db *repository.DBs, grpcCli exerpb.ExerciseServiceClient) ServiceIface {
	return &service{Db: Db, GClient: grpcCli}
}

func (s *service) CreatePlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) error {
	var exerciseIDs []string
	PlanName := strings.ToLower(planName)

	for _, exerciseName := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: exerciseName, UserId: userId})
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

func (s *service) GetAllPlansSer(ctx context.Context, userId string) (int, *[]models.Plan2, error) {

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
			exerciseName, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: id, UserId: userId})
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

func (s *service) GetPlanByNameSer(ctx context.Context, userId string, planName string) (string, string, *[]string, error) {

	planId, err := s.Db.ReturnsPlanId(ctx, userId, planName)
	if err != nil {
		return "", "", nil, err
	}

	exerciseIds, err := s.Db.GetAllExercisesByPlanID(ctx, planId)
	if err != nil {
		return "", "", nil, err
	}

	var allExercises []string
	for _, exerciseId := range *exerciseIds {

		r, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: exerciseId, UserId: userId})
		if err != nil {
			return "", "", nil, err
		}

		allExercises = append(allExercises, r.ExerciseName)
	}

	return planId, planName, &allExercises, nil
}

func (s *service) AddExercisesToPlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error) {

	planId, err := s.Db.ReturnsPlanId(ctx, userId, planName)
	if err != nil {
		return nil, err
	}

	var exerciseIds []string
	for _, eachName := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: eachName, UserId: userId})
		if err != nil {
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

func (s *service) DeleteExerciseFromPlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error) {

	planId, err := s.Db.ReturnsPlanId(ctx, userId, planName)
	if err != nil {
		return nil, err
	}

	var exerciseIds []string

	for _, v := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: v, UserId: userId})
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

func (s *service) DeletePlanSer(ctx context.Context, userId string, planName string) error {

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


func (s *service) GetEmptyPlanId(ctx context.Context, userId string) (string, error) {

	var PlanId string

	PlanId, err := s.Db.GetUserEmptyPlanIdR(ctx, userId)
	if err != nil {
		return "", err
	}
	
	if PlanId == "" {
		PlanId, err := s.Db.ReturnsPlanId(ctx, userId, "empty")
		if err != nil {
			return "", err
		}
		
		err = s.Db.SetUserEmptyPlanIdR(ctx, userId, PlanId)
		if err != nil {
			return "", err
		}
	}

	return PlanId, nil
}

func MakeRespForAddingNewExer(ctx context.Context, userId string, planId string, planName string, s *service) (*models.Plan2, error) {

	var resp models.Plan2
	var allExercises []string

	exerciseIds, err := s.Db.GetAllExercisesByPlanID(ctx, planId)
	if err != nil {
		return &resp, err
	}

	for _, exerciseId := range *exerciseIds {

		r, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: exerciseId, UserId: userId})
		if err != nil {
			return &resp, fmt.Errorf("error getting data from exercise grpc server : %w", err)
		}

		allExercises = append(allExercises, r.ExerciseName)
	}

	resp.PlanName = planName
	resp.Exercises = allExercises

	return &resp, nil
}

func (s *service) GetHealth(ctx context.Context) (*time.Duration, *time.Duration) {

	pgRespTime := s.Db.GetPostgresRespTime(ctx)
	redisRespTime := s.Db.GetRedisRespTime(ctx)

	return pgRespTime, redisRespTime
}
