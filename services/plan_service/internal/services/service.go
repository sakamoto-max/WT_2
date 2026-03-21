package services

import (
	"context"
	"fmt"
	customerrors "plan_service/internal/custom_errors"
	"plan_service/internal/models"
	"plan_service/internal/repository"
	"time"

	// "plan_service/internal/user"

	// "strconv"
	"strings"
	exerpb "workout-tracker/proto/shared/exercise"
)

type Service struct {
	Db      *repository.DBs
	GClient exerpb.ExerciseServiceClient
}

func NewService(Db *repository.DBs, grpcCli exerpb.ExerciseServiceClient) *Service {
	return &Service{Db: Db, GClient: grpcCli}
}

func (s *Service) CreatePlan(ctx context.Context, userId int, planName string, exerciseNames *[]string) error {
	var exerciseIDs []int
	// var resp models.Plan2Resp
	// lower case the plan_name
	// replace " " with _  (TODO)
	PlanName := strings.ToLower(planName)
	// check if plan already exists
	exists, err := s.Db.PlanExists(ctx, userId, PlanName)
	if err != nil {
		return err
	}

	if exists {
		return customerrors.ErrPlanAlreadyExists
	}

	for _, exerciseName := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: exerciseName})
		if err != nil {
			return fmt.Errorf("error getting data from execise server : %w", err)
		}

		if !r.Exists {
			return fmt.Errorf("exercise %v doesnot exits, please create it : %w\n", exerciseName, err)
		}

		exerciseIDs = append(exerciseIDs, int(r.ExerciseId))
	}

	err = s.Db.CreatePlan(ctx, userId, PlanName, exerciseIDs)
	if err != nil {
		return err
	}

	return nil
}
func (s *Service) GetAllPlansSer(ctx context.Context, userId int) (*models.AllPlansResp, error) {
	// get all plan Ids of the user

	var resp models.AllPlansResp

	var allPlans []models.Plan2

	planNamesWithIds, err := s.Db.GetAllUserPlans(ctx, userId)
	if err != nil {
		return &resp, err
	}

	for _, eachPlan := range *planNamesWithIds {

		var plan models.Plan2
		var exerciseNames []string

		exeriseIDs, err := s.Db.GetAllExercisesByPlanID(ctx, eachPlan.Id)
		if err != nil {
			return &resp, err
		}
		// get the name for each exerciseid from redis

		for _, v := range *exeriseIDs {
			// id := strconv.Itoa(v)
			exerciseName, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: int64(v)})
			// name, err := s.Db.GetExerciseNameByID(ctx, id)
			if err != nil {
				return &resp, err
			}

			exerciseNames = append(exerciseNames, exerciseName.ExerciseName)
		}

		plan.PlanName = eachPlan.PlanName
		plan.Exercises = exerciseNames
		allPlans = append(allPlans, plan)
	}

	numberOfPlans := len(*planNamesWithIds)
	resp.NumberOfPlans = numberOfPlans
	resp.Plans = allPlans
	return &resp, nil
}
func (s *Service) GetPlanByNameSer(ctx context.Context, userId int, planName string) (*models.Plan2, error) {

	var resp models.Plan2
	var allExercises []string

	exists, planId, err := s.Db.PlanExistsReturnsId(ctx, userId, planName)
	if err != nil {
		return &resp, err
	}

	if !exists {
		return &resp, customerrors.ErrPlanNameDoesNotExists
	}

	exerciseIds, err := s.Db.GetAllExercisesByPlanID(ctx, planId)
	if err != nil {
		return &resp, err
	}

	for _, exerciseId := range *exerciseIds {

		r, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: int64(exerciseId)})
		if err != nil {
			return &resp, fmt.Errorf("error getting data from exercise grpc server : %w", err)
		}

		allExercises = append(allExercises, r.ExerciseName)
	}

	resp.PlanName = planName
	resp.Exercises = allExercises

	return &resp, nil
}
func (s *Service) AddExercisesToPlan(ctx context.Context, userId int, planName string, exerciseNames *[]string) (*models.Plan2, error) {

	var exerciseIds []int
	var resp *models.Plan2

	// check if plan exists
	exists, planId, err := s.Db.PlanExistsReturnsId(ctx, userId, planName)
	if err != nil {
		//
		return resp, err
	}

	if !exists {
		return resp, customerrors.ErrPlanNameDoesNotExists
	}

	// check if exercise exists
	// get all the ids of exercises from grpc
	for _, v := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: v})
		if err != nil {
			return resp, fmt.Errorf("error getting data from exercise server : %w", err)
		}

		if !r.Exists {
			return resp, fmt.Errorf("exercise %v does not exist", v)
		}

		exerciseIds = append(exerciseIds, int(r.ExerciseId))
	}

	err = s.Db.AddExercisesToPlan(ctx, planId, &exerciseIds)
	if err != nil {
		return resp, err
	}

	resp, err = MakeRespForAddingNewExer(ctx, planId, planName, s)
	if err != nil {
		return resp, err
	}

	return resp, nil

}
func (s *Service) DeleteExerciseFromPlan(ctx context.Context, userId int, planName string, exerciseNames *[]string) (*models.Plan2, error) {

	// get plan
	var exerciseIds []int
	var resp *models.Plan2

	exists, planId, err := s.Db.PlanExistsReturnsId(ctx, userId, planName)
	if err != nil {
		//
		return resp, err
	}

	if !exists {
		return resp, customerrors.ErrPlanNameDoesNotExists
	}

	for _, v := range *exerciseNames {

		r, err := s.GClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: v})
		if err != nil {
			return resp, fmt.Errorf("error getting data from exercise server : %w", err)
		}

		if !r.Exists {
			return resp, fmt.Errorf("exercise %v does not exist", v)
		}

		exerciseIds = append(exerciseIds, int(r.ExerciseId))
	}

	err = s.Db.DeleteExerciseFromPlan(ctx, planId, &exerciseIds)
	if err != nil {
		return resp, err
	}

	resp, err = MakeRespForAddingNewExer(ctx, planId, planName, s)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
func (s *Service) DeletePlanSer(ctx context.Context, userId int, planName string) error {
	// check if plan exists -> gets plan id

	exists, planId, err := s.Db.PlanExistsReturnsId(ctx, userId, planName)
	if err != nil {
		return err
	}

	if !exists {
		return customerrors.ErrPlanNameDoesNotExists
	}

	err = s.Db.DeletePlan(ctx, userId, planId)
	if err != nil {
		return err
	}

	return nil
}
func (s *Service) PlanExistsReturnId(ctx context.Context, userID int, planName string) (bool, int, error) {
	return s.Db.PlanExistsReturnsId(ctx, userID, planName)
}
func (s *Service) PlanExistsReturnPlan(ctx context.Context, userId int, planName string) (bool, int, *[]int, error) {
	return s.Db.PlanExistsReturnPlan(ctx, userId, planName)
}

func (s *Service) GetEmptyPlanId(ctx context.Context, userId int) (int, error) {
	return s.Db.GetEmptyPlanID(ctx, userId)
}

func (s *Service) CreateEmptyPlan(ctx context.Context, userId int) error {
	return s.Db.CreateEmptyPlan(ctx, userId)
}

func MakeRespForAddingNewExer(ctx context.Context, planId int, planName string, s *Service) (*models.Plan2, error) {

	var resp models.Plan2
	var allExercises []string

	exerciseIds, err := s.Db.GetAllExercisesByPlanID(ctx, planId)
	if err != nil {
		return &resp, err
	}

	for _, exerciseId := range *exerciseIds {

		r, err := s.GClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: int64(exerciseId)})
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
