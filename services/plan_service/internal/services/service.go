package services

import (
	"context"
	"fmt"
	"plan_service/internal/models"
	"plan_service/internal/repository"
	"strings"
	"time"

	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
)

type ServiceIface interface {
	CreatePlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) error
	GetPlans(ctx context.Context, userId string) (int, *[]models.Plan2, error)
	GetPlan(ctx context.Context, userId string, planName string) (string, string, *[]string, error)
	AddExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error)
	DeleteExerciseFromPlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error)
	DeletePlan(ctx context.Context, userId string, planName string) error
	GetEmptyPlanId(ctx context.Context, userId string) (string, error)
	GetHealth(ctx context.Context) (*time.Duration, *time.Duration)
}

type service struct {
	Db      repository.RepoIFace
	GClient exerpb.ExerciseServiceClient
}

func NewService(Db repository.RepoIFace, grpcCli exerpb.ExerciseServiceClient) ServiceIface {
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

func (s *service) GetPlans(ctx context.Context, userId string) (int, *[]models.Plan2, error) {

	var allPlans []models.Plan2

	planNamesWithIds, err := s.Db.GetPlans(ctx, userId)
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

func (s *service) GetPlan(ctx context.Context, userId string, planName string) (string, string, *[]string, error) {

	PlanId, ExerciseIds, err := s.Db.GetUserPlan(ctx, userId, planName)
	if err != nil {
		return "", "", nil, err
	}

	if PlanId == "" {
		PlanId, err = s.Db.ReturnsPlanId(ctx, userId, planName)
		if err != nil {
			return "", "", nil, err
		}

		ExerciseIds, err = s.Db.GetAllExercisesByPlanID(ctx, PlanId)
		if err != nil {
			return "", "", nil, err
		}

		err = s.Db.SetUserPlan(ctx, userId, planName, PlanId, ExerciseIds)
		if err != nil {
			return "", "", nil, err
		}
	}

	var allExercises []string
	for _, exerciseId := range *ExerciseIds {

		r, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: exerciseId, UserId: userId})
		if err != nil {
			return "", "", nil, err
		}

		allExercises = append(allExercises, r.ExerciseName)
	}

	return PlanId, planName, &allExercises, nil
}

func (s *service) AddExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error) {

	PlanId, err := s.Db.GetUserPlanId(ctx, userId, planName)
	if err != nil {
		return nil, err
	}

	if PlanId == "" {
		PlanId, err = s.Db.ReturnsPlanId(ctx, userId, planName)
		if err != nil {
			return nil, err
		}

		err = s.Db.SetUserPlanId(ctx, userId, planName, PlanId)
		if err != nil {
			return nil, err
		}
	}

	var exerciseIds []string
	for _, eachName := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: eachName, UserId: userId})
		if err != nil {
			return nil, err
		}

		exerciseIds = append(exerciseIds, r.ExerciseId)
	}

	err = s.Db.AddExercisesToPlan(ctx, PlanId, &exerciseIds)
	if err != nil {
		return nil, err
	}

	resp, err := makeRespForAddingNewExer(ctx, userId, PlanId, planName, s)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) DeleteExerciseFromPlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error) {

	PlanId, err := s.Db.GetUserPlanId(ctx, userId, planName)
	if err != nil {
		return nil, err
	}

	if PlanId == "" {
		PlanId, err = s.Db.ReturnsPlanId(ctx, userId, planName)
		if err != nil {
			return nil, err
		}

		err = s.Db.SetUserPlanId(ctx, userId, planName, PlanId)
		if err != nil {
			return nil, err
		}
	}

	var exerciseIds []string

	for _, v := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: v, UserId: userId})
		if err != nil {
			return nil, fmt.Errorf("error getting data from exercise server : %w", err)
		}

		exerciseIds = append(exerciseIds, r.ExerciseId)
	}

	err = s.Db.DeleteExerciseFromPlan(ctx, PlanId, &exerciseIds)
	if err != nil {
		return nil, err
	}

	resp, err := makeRespForAddingNewExer(ctx, userId, PlanId, planName, s)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) DeletePlan(ctx context.Context, userId string, planName string) error {

	
	PlanId, err := s.Db.GetUserPlanId(ctx, userId, planName)
	if err != nil {
		return err
	}
	
	if PlanId == "" {
		PlanId, err = s.Db.ReturnsPlanId(ctx, userId, planName)
		if err != nil {
			return err
		}
	}

	err = s.Db.DeletePlan(ctx, userId, PlanId)
	if err != nil {
		return err
	}

	err = s.Db.DelUserPlanId(ctx, userId, planName)
	if err != nil {
		return err
	}

	err = s.Db.DelUserPlan(ctx, userId, planName)
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

func makeRespForAddingNewExer(ctx context.Context, userId string, planId string, planName string, s *service) (*models.Plan2, error) {

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
