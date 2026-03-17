package services

import (
	"context"
	// "fmt"
	// grpcclient "plan_service/grpc_client"
	"plan_service/internal/repository"

	// pb "workout-tracker/proto/shared/plan"
	// epb "workout-tracker/proto/shared/exercise"
	pb "workout-tracker/proto/shared/plan"
)

type PlanService struct {
	dB *repository.DBs
}

func NewPlanService(db *repository.DBs) *PlanService {
	return &PlanService{dB: db}
}

func (p *PlanService) PlanExistsReturnId(ctx context.Context, userID int, planName string) (bool, int, error) {
	return p.dB.PlanExistsReturnsId(ctx, userID, planName)
}

func (p *PlanService) PlanExistsReturnPlan(ctx context.Context, userId int, planName string) (*pb.PlanExistsReturnPlanResp, error) {
	var resp pb.PlanExistsReturnPlanResp
	exists, planId, exerciseIds, err := p.dB.PlanExistsReturnPlan(ctx, userId, planName)
	if err != nil {
		return &resp, err
	}

	resp.Exists = exists
	resp.PlanId = int64(planId)
	for _, v := range *exerciseIds {
		resp.ExerciseIds = append(resp.ExerciseIds, int64(v))
	}

	return &resp, nil
}

func (p *PlanService) GetEmptyPlanId(ctx context.Context, userId int) (*pb.EmptyPlanIdResp, error) {
	var resp pb.EmptyPlanIdResp
	emptyPlanId, err := p.dB.GetEmptyPlanID(ctx, userId)
	if err != nil {
		return &resp, err
	}

	resp.EmptyPlanId = int64(emptyPlanId)

	return &resp, nil
}

func (p *PlanService) CreateEmptyPlan(ctx context.Context, userId int) (*pb.CreateEmptyPlanResp, error) {
	var resp pb.CreateEmptyPlanResp
	err := p.dB.CreateEmptyPlan(ctx, userId)
	if err != nil {
		return &resp, err
	}

	resp.Message = "empty plan created successfully"

	return &resp, err
}
