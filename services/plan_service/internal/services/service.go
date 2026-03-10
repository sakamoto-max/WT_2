package services

import (
	"context"
	"fmt"
	customerrors "plan_service/internal/custom_errors"
	"plan_service/internal/models"
	"plan_service/internal/repository"
	"strconv"
	"strings"
)

type Service struct {
	Db      *repository.DBs
}

func NewService(Db *repository.DBs) *Service {
	return &Service{Db: Db}
}

func (s *Service) CreatePlanSer(ctx context.Context, userId int, plan *models.Plan2) (*models.Plan2Resp, error) {
	var resp models.Plan2Resp
	// lower case the plan_name
	// replace " " with _  (TODO)
	planName := strings.ToLower(plan.PlanName)
	// check if plan already exists
	exists, err := s.Db.PlanExists(ctx, userId, planName)
	if err != nil {
		return &resp, err
	}

	if exists {
		return &resp, customerrors.ErrPlanAlreadyExists
	}

	for _, exerciseName := range plan.Exercises {
		exists, err := s.Db.ExerciseExistsInMain(ctx, exerciseName)
		if err != nil {
			return &resp, err
		}

		if !exists {
			return &resp, fmt.Errorf("exercise %v doesnot exits, please create it : %w\n", exerciseName, err)
		}
	}

	err = s.Db.CreatePlan(ctx, userId, plan)
	if err != nil {
		return &resp, err
	}

	resp.PlanName = plan.PlanName
	resp.Exercises = plan.Exercises
	resp.Message = fmt.Sprintf("%v created successfully", plan.PlanName)

	return &resp, nil
}


func (s *Service) GetAllPlansSer(ctx context.Context, userId int) (*models.AllPlansResp, error) {
	// get all plan Ids of the user

	var resp models.AllPlansResp

	var allPlans []models.Plan2

	planIds, err := s.Db.GetAllUsersPlanIds(ctx, userId)
	if err != nil {
		return &resp, err
	}

	for _, planId := range *planIds {

		var plan models.Plan2
		var exerciseNames []string

		// get planName
		planName, err := s.Db.GetPlanNameByID(ctx, planId)
		if err != nil {
			return &resp, err
		}

		// get all the exercise ids of this plan
		exeriseIDs, err := s.Db.GetAllExercisesByPlanID(ctx, planId)
		if err != nil {
			return &resp, err
		}
		// get the name for each exerciseid from redis

		for _, v := range *exeriseIDs {
			id := strconv.Itoa(v)
			name, err := s.Db.GetExerciseNameByID(ctx, id)
			if err != nil {
				return &resp, err
			}

			exerciseNames = append(exerciseNames, name)
		}

		plan.PlanName = planName
		plan.Exercises = exerciseNames
		allPlans = append(allPlans, plan)
	}

	numberOfPlans := len(*planIds)
	resp.NumberOfPlans = numberOfPlans
	resp.Plans = allPlans
	return &resp, nil
}

func (s *Service) GetPlanByNameSer(ctx context.Context, userId int, planName string) (*models.Plan2, error) {

	var resp models.Plan2
	var allExercises []string

	loaded, err := s.Db.UserPlanNamesLoaded(ctx, userId)
	if err != nil {
		return &resp, err
	}

	if !loaded {
		err := s.Db.SetAllUserPlanNames(ctx, userId)
		if err != nil {
			return &resp, err
		}
	}

	planID, err := s.Db.GetPlanIdFromRedis(ctx, userId, planName)
	if err != nil {
		return &resp, err
	}

	// get the exercises in the plan

	exerciseIds, err := s.Db.GetAllExercisesByPlanID(ctx, planID)
	if err != nil {
		return &resp, err
	}

	for _, exerciseId := range *exerciseIds {

		id := strconv.Itoa(exerciseId)

		exerciseName, err := s.Db.GetExerciseNameByID(ctx, id)
		if err != nil {
			return &resp, err
		}

		allExercises = append(allExercises, exerciseName)
	}

	resp.PlanName = planName
	resp.Exercises = allExercises

	return &resp, nil
}
