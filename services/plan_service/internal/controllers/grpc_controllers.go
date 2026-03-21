package controllers

import (
	"context"
	"fmt"
	"plan_service/internal/services"
	planpb "workout-tracker/proto/shared/plan"

	"google.golang.org/protobuf/types/known/durationpb"
)

type PlanController struct {
	service *services.Service
	planpb.UnimplementedPlanServiceServer
}

func NewPlanController(service *services.Service) *PlanController {
	return &PlanController{service: service}
}
func (p *PlanController) CreatePlan(ctx context.Context, in *planpb.CreatePlanReq) (*planpb.CreatePlanResp, error) {
	resp := planpb.CreatePlanResp{}

	err := p.service.CreatePlan(ctx, int(in.UserId), in.PlanName, &in.ExerciseNames)
	if err != nil {
		return &resp, err
	}

	resp.PlanName = in.PlanName
	resp.ExerciseNames = in.ExerciseNames
	resp.Message = fmt.Sprintf("%v created successfully", in.PlanName)

	return &resp, nil
}
func (p *PlanController) CheckHealth(ctx context.Context, in *planpb.CheckHealthReq) (*planpb.CheckHealthResp, error) {
	resp := planpb.CheckHealthResp{}

	// ping the dbs

	resp.Message = "every thing is all right"

	return &resp, nil
}
func (p *PlanController) PlanExistsReturnId(ctx context.Context, in *planpb.SendPlanName) (*planpb.PlanExistsResp, error) {
	resp := planpb.PlanExistsResp{}
	exits, planId, err := p.service.PlanExistsReturnId(ctx, int(in.UserId), in.PlanName)
	if err != nil {
		return &resp, err
	}

	resp.Exists = exits
	resp.PlanId = int64(planId)

	return &resp, nil
}
func (p *PlanController) GetEmptyPlanId(ctx context.Context, in *planpb.SendUserID) (*planpb.EmptyPlanIdResp, error) {
	resp := planpb.EmptyPlanIdResp{}

	emptyPlanId, err := p.service.GetEmptyPlanId(ctx, int(in.UserId))

	if err != nil {
		return &resp, err
	}

	resp.EmptyPlanId = int64(emptyPlanId)

	return &resp, nil
}
func (p *PlanController) PlanExistsReturnPlan(ctx context.Context, in *planpb.SendPlanName) (*planpb.PlanExistsReturnPlanResp, error) {
	resp := planpb.PlanExistsReturnPlanResp{}

	exists, planId, exerciseIds, err := p.service.PlanExistsReturnPlan(ctx, int(in.UserId), in.PlanName)
	if err != nil {
		return &resp, err
	}

	resp.Exists = exists
	resp.PlanId = int64(planId)
	for _, eachExerciseId := range *exerciseIds {
		resp.ExerciseIds = append(resp.ExerciseIds, int64(eachExerciseId))
	}

	return &resp, nil
}
func (p *PlanController) CreateEmptyPlan(ctx context.Context, in *planpb.SendUserID) (*planpb.CreateEmptyPlanResp, error) {
	resp := planpb.CreateEmptyPlanResp{}

	err := p.service.CreateEmptyPlan(ctx, int(in.UserId))
	if err != nil {
		return &resp, err
	}

	resp.Message = "empty plan created successfully"

	return &resp, nil

}
func (p *PlanController) GetAllPlans(ctx context.Context, in *planpb.GetAllPlansReq) (*planpb.GetAllPlansResp, error) {
	resp := planpb.GetAllPlansResp{}
	r, err := p.service.GetAllPlansSer(ctx, int(in.UserId))
	if err != nil {
		return &resp, err
	}

	for _, v := range r.Plans {
		eachPlan := planpb.PlanResp{}

		eachPlan.PlanName = v.PlanName
		eachPlan.ExerciseNames = v.Exercises

		resp.AllPlans = append(resp.AllPlans, &eachPlan)
	}
	resp.NumberOfPlans = int64(r.NumberOfPlans)

	return &resp, nil
}
func (p *PlanController) GetPlanByName(ctx context.Context, in *planpb.GetPlanByNameReq) (*planpb.PlanResp, error) {
	resp := planpb.PlanResp{}
	r, err := p.service.GetPlanByNameSer(ctx, int(in.UserId), in.PlanName)
	if err != nil {
		return &resp, err
	}

	resp.ExerciseNames = r.Exercises
	resp.PlanName = r.PlanName

	return &resp, nil
}
func (p *PlanController) AddExercisesToPlan(ctx context.Context, in *planpb.PlanReq) (*planpb.PlanResp, error) {
	resp := planpb.PlanResp{}

	r, err := p.service.AddExercisesToPlan(ctx, int(in.UserId), in.PlanName, &in.ExerciseNames)
	if err != nil {
		return &resp, err
	}

	resp.PlanName = r.PlanName
	resp.ExerciseNames = r.Exercises

	return &resp, nil
}
func (p *PlanController) DeleteExercisesFromPlan(ctx context.Context, in *planpb.PlanReq) (*planpb.PlanResp, error) {
	resp := planpb.PlanResp{}
	r, err := p.service.DeleteExerciseFromPlan(ctx, int(in.UserId), in.PlanName, &in.ExerciseNames)
	if err != nil {
		return &resp, err
	}

	resp.PlanName = r.PlanName
	resp.ExerciseNames = r.Exercises

	return &resp, nil

}
func (p *PlanController) DeletePlan(ctx context.Context, in *planpb.DeletePlanReq) (*planpb.DeletePlanResp, error) {
	resp := planpb.DeletePlanResp{}
	err := p.service.DeletePlanSer(ctx, int(in.UserId), in.PlanName)
	if err != nil {
		return &resp, err
	}

	return &resp, nil
}

func (a *PlanController) PING(ctx context.Context, in *planpb.PingPlanReq) (*planpb.PingPlanResp, error) {
	r := planpb.PingPlanResp{}

	return &r, nil
}

func (a *PlanController) GetHealth(ctx context.Context, in *planpb.GetHealthReq) (*planpb.GetHealthResp, error) {

	resp := planpb.GetHealthResp{}

	pgRespTime, redisRespTime := a.service.GetHealth(ctx)

	resp.PostgresRespTime = durationpb.New(*pgRespTime)
	resp.RedisRespTime = durationpb.New(*redisRespTime)

	return &resp, nil
}
