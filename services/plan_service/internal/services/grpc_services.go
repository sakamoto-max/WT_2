package services

import (
	"context"
	"plan_service/internal/repository"
	// pb "workout-tracker/proto/shared/plan"
	pb "github.com/sakamoto-max/wt_2-proto/shared/plan"

)

type PlanService struct {
	dB *repository.DBs
}

func NewPlanService(db *repository.DBs) *PlanService {
	return &PlanService{dB: db}
}

func (p *PlanService) CreateEmptyPlan(ctx context.Context, userId string) (*pb.CreateEmptyPlanResp, error) {

	err := p.dB.CreateEmptyPlan(ctx, userId)
	if err != nil {
		return nil, err
	}

	resp := pb.CreateEmptyPlanResp{
		Message: "empty plan created successfully",
	}

	return &resp, err
}
