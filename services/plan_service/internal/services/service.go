package services

import (
	"context"
	"errors"
	"fmt"
	"plan_service/internal/domain"
	"plan_service/internal/repository"
	"strings"

	"github.com/sakamoto-max/wt_2_pkg/myerrs"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
)

func (s *Service) CreatePlan(ctx context.Context, in *planpb.CreatePlanReq) (*planpb.CreatePlanResp, error) {

	var exerciseIDs []string

	for _, exerciseName := range in.ExerciseNames {

		r, err := s.gClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: exerciseName, UserId: in.UserId})
		if err != nil {
			return nil, err
		}

		exerciseIDs = append(exerciseIDs, r.ExerciseId)
	}

	in.PlanName = strings.ToLower(in.PlanName)

	mappedData := domain.ToCreatePlan(in)

	mappedData.ExerciseIds = &exerciseIDs

	err := s.pg.PlanCommandRepo.CreatePlan(ctx, mappedData)
	if err != nil {
		return nil, err
	}

	return &planpb.CreatePlanResp{
		PlanName:      in.PlanName,
		ExerciseNames: in.ExerciseNames,
		Message:       fmt.Sprintf("%v created successfully", in.PlanName),
	}, nil
}

func (s *Service) GetAllPlans(ctx context.Context, in *planpb.GetAllPlansReq) (*planpb.GetAllPlansResp, error) {

	var allPlans []domain.Plan

	planNamesWithIds, err := s.pg.PlanQueryRepo.GetAllPlanNamesWithIds(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	for _, eachPlan := range *planNamesWithIds {

		var plan domain.Plan
		var exerciseNames []string

		exeriseIDs, err := s.pg.PlanQueryRepo.GetAllExercisesByPlanID(ctx, eachPlan.PlanId)
		if err != nil {
			return nil, err
		}

		for _, id := range *exeriseIDs {
			exerciseName, err := s.gClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: id, UserId: in.UserId})
			if err != nil {
				return nil, err
			}

			exerciseNames = append(exerciseNames, exerciseName.ExerciseName)
		}

		plan.PlanName = eachPlan.PlanName
		plan.ExerciseNames = &exerciseNames
		allPlans = append(allPlans, plan)
	}

	resp := planpb.GetAllPlansResp{}

	for _, eachPlan := range allPlans {
		plan := planpb.PlanResp{
			ExerciseNames: *eachPlan.ExerciseNames,
			PlanName:      eachPlan.PlanName,
		}

		resp.AllPlans = append(resp.AllPlans, &plan)
	}

	resp.NumberOfPlans = int64(len(resp.AllPlans))

	return &resp, nil
}

func (s *Service) GetPlanByName(ctx context.Context, in *planpb.GetPlanByNameReq) (*planpb.GetPlanByNameResp, error) {

	fmt.Println("get plan by name called")

	PlanId, ExerciseIds, err := s.cache.UserPlan.GetUserPlan(ctx, domain.ToGetPlan(in))
	if err != nil || PlanId == "" {

		PlanId, err = s.pg.PlanQueryRepo.GetPlanId(ctx, domain.ToGetPlan(in))
		if err != nil {
			return nil, err
		}

		ExerciseIds, err = s.pg.PlanQueryRepo.GetAllExercisesByPlanID(ctx, PlanId)
		if err != nil {
			return nil, err
		}

		plan := domain.Plan{
			UserId:      in.UserId,
			PlanId:      PlanId,
			PlanName:    in.PlanName,
			ExerciseIds: ExerciseIds,
		}

		s.cache.UserPlan.SetUserPlan(ctx, plan)
	}

	var allExerciseNames []string
	for _, exerciseId := range *ExerciseIds {

		r, err := s.gClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: exerciseId, UserId: in.UserId})
		if err != nil {
			return nil, err
		}

		allExerciseNames = append(allExerciseNames, r.ExerciseName)
	}

	fmt.Println("plan in plan", planpb.GetPlanByNameResp{
		PlanName:      in.PlanName,
		PlanId:        PlanId,
		ExerciseNames: allExerciseNames,
	})

	return &planpb.GetPlanByNameResp{
		PlanName:      in.PlanName,
		PlanId:        PlanId,
		ExerciseNames: allExerciseNames,
	}, nil
}

func (s *Service) AddExercisesToPlan(ctx context.Context, in *planpb.PlanReq) (*planpb.PlanResp, error) {

	PlanId, err := s.cache.PlanId.GetUserPlanId(ctx, domain.ToGetPlan(in))
	if err != nil || PlanId == "" {
		PlanId, err = s.pg.PlanQueryRepo.GetPlanId(ctx, domain.ToGetPlan(in))
		if err != nil {
			return nil, err
		}

		s.cache.PlanId.SetUserPlanId(ctx, domain.ToGetPlan(in), PlanId)
	}

	var exerciseIds []string
	for _, eachName := range in.ExerciseNames {

		r, err := s.gClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: eachName, UserId: in.UserId})
		if err != nil {
			return nil, err
		}

		exerciseIds = append(exerciseIds, r.ExerciseId)
	}

	err = s.pg.PlanExericseRepo.AddExercisesToPlan(ctx, PlanId, &exerciseIds)
	if err != nil {
		var dbErr *repository.DbErr
		if errors.As(err, &dbErr) {
			exerId := dbErr.GetExerciseId()
			resp, err := s.gClient.GetExerciseName(ctx, &exerpb.SendExerciseID{UserId: in.UserId, ExerciseId: exerId})
			if err != nil {
				return nil, myerrs.InternalServerErrMaker(err)
			}

			return nil, myerrs.AlreadyExitsErrMaker(fmt.Sprintf("exercise %s", resp.ExerciseName))
		}

		return nil, err
	}

	allExerciseIds, err := s.pg.PlanQueryRepo.GetAllExercisesByPlanID(ctx, PlanId)
	if err != nil {
		return nil, err
	}

	var allExerciseNames []string

	for _, exerciseId := range *allExerciseIds {

		r, err := s.gClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: exerciseId, UserId: in.UserId})
		if err != nil {
			return nil, fmt.Errorf("error getting data from exercise grpc server : %w", err)
		}

		allExerciseNames = append(allExerciseNames, r.ExerciseName)
	}

	if err = s.cache.UserPlan.DelUserPlan(ctx, domain.GetPlan{UserId: in.UserId, PlanName: in.PlanName}); err != nil {
		return nil, err
	}

	return &planpb.PlanResp{
		PlanName:      in.PlanName,
		ExerciseNames: allExerciseNames,
	}, nil
}

func (s *Service) GetEmptyPlanId(ctx context.Context, in *planpb.SendUserID) (*planpb.EmptyPlanIdResp, error) {

	PlanId, err := s.cache.EmptyPlan.GetUserEmptyPlanId(ctx, in.UserId)
	if err != nil || PlanId == "" {

		PlanId, err = s.pg.PlanQueryRepo.GetEmptyPlanId(ctx, in.UserId)
		if err != nil {
			return nil, err
		}

		s.cache.EmptyPlan.SetUserEmptyPlanId(ctx, in.UserId, PlanId)
	}

	return &planpb.EmptyPlanIdResp{
		EmptyPlanId: PlanId,
	}, nil
}

func (s *Service) DeletePlan(ctx context.Context, in *planpb.DeletePlanReq) (*planpb.DeletePlanResp, error) {
	PlanId, err := s.cache.PlanId.GetUserPlanId(ctx, domain.ToGetPlan(in))
	if err != nil || PlanId == "" {
		PlanId, err = s.pg.PlanQueryRepo.GetPlanId(ctx, domain.ToGetPlan(in))
		if err != nil {
			return nil, err
		}
	}

	err = s.pg.PlanCommandRepo.DeletePlan(ctx, in.UserId, PlanId)
	if err != nil {
		return nil, err
	}

	s.cache.PlanId.DelUserPlanId(ctx, domain.ToGetPlan(in))

	s.cache.UserPlan.DelUserPlan(ctx, domain.ToGetPlan(in))

	return &planpb.DeletePlanResp{}, nil
}

func (s *Service) DeleteExercisesFromPlan(ctx context.Context, in *planpb.PlanReq) (*planpb.PlanResp, error) {

	PlanId, err := s.cache.PlanId.GetUserPlanId(ctx, domain.ToGetPlan(in))
	if err != nil || PlanId == "" {
		PlanId, err = s.pg.PlanQueryRepo.GetPlanId(ctx, domain.ToGetPlan(in))
		if err != nil {
			return nil, err
		}

		s.cache.PlanId.SetUserPlanId(ctx, domain.ToGetPlan(in), PlanId)
	}

	var exerciseIds []string

	for _, v := range in.ExerciseNames {

		r, err := s.gClient.ExerciseExistsReturnId(ctx, &exerpb.SendExerciseName{ExerciseName: v, UserId: in.UserId})
		if err != nil {
			return nil, fmt.Errorf("error getting data from exercise server : %w", err)
		}

		exerciseIds = append(exerciseIds, r.ExerciseId)
	}

	err = s.pg.PlanExericseRepo.RemoveExerciseFromPlan(ctx, PlanId, &exerciseIds)
	if err != nil {
		return nil, err
	}

	allExerciseIds, err := s.pg.PlanQueryRepo.GetAllExercisesByPlanID(ctx, PlanId)
	if err != nil {
		return nil, err
	}

	var allExerciseNames []string

	for _, exerciseId := range *allExerciseIds {

		r, err := s.gClient.GetExerciseName(ctx, &exerpb.SendExerciseID{ExerciseId: exerciseId, UserId: in.UserId})
		if err != nil {
			return nil, fmt.Errorf("error getting data from exercise grpc server : %w", err)
		}

		allExerciseNames = append(allExerciseNames, r.ExerciseName)
	}

	if err = s.cache.UserPlan.DelUserPlan(ctx, domain.GetPlan{UserId: in.UserId, PlanName: in.PlanName}); err != nil {
		return nil, err
	}

	return &planpb.PlanResp{
		PlanName:      in.PlanName,
		ExerciseNames: allExerciseNames,
	}, nil
}
