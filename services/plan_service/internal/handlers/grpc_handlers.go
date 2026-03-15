package handlers

import (
	"context"
	"plan_service/internal/services"
	pb "workout-tracker/proto/shared/plan"
)

type PlanHandler struct {
	service *services.PlanService
	pb.UnimplementedPlanServiceServer
}

func NewgRPCHandler(s *services.PlanService) *PlanHandler {
	return &PlanHandler{
		service: s,
	}
}

func (p *PlanHandler) PlanExistsReturnId(ctx context.Context, req *pb.SendPlanName) (*pb.PlanExistsResp, error) {
	var resp pb.PlanExistsResp

	exists, planId, err := p.service.PlanExistsReturnId(ctx, int(req.UserId), req.PlanName)
	if err != nil{
		return &resp, err
	}
	resp.Exists = exists
	resp.PlanId = int64(planId)

	return &resp, nil
}

func (p *PlanHandler) PlanExistsReturnPlan(ctx context.Context, in *pb.SendPlanName) (*pb.PlanExistsReturnPlanResp, error) {
	return p.service.PlanExistsReturnPlan(ctx, int(in.UserId), in.PlanName)
}

func (p *PlanHandler) GetEmptyPlanId(ctx context.Context, in *pb.SendUserID) (*pb.EmptyPlanIdResp, error) {
	return p.service.GetEmptyPlanId(ctx, int(in.UserId))
}

func (p *PlanHandler) CreateEmptyPlan(ctx context.Context, in *pb.SendUserID) (*pb.CreateEmptyPlanResp, error) {
	return p.service.CreateEmptyPlan(ctx, int(in.UserId))
}
