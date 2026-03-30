package controllers

import (
	"context"
	"fmt"
	"plan_service/internal/services"
	planpb "workout-tracker/proto/shared/plan"
	// "google.golang.org/protobuf/types/known/durationpb"
)

type PlanController struct {
	service *services.Service
	planpb.UnimplementedPlanServiceServer
}

func NewPlanController(service *services.Service) *PlanController {
	return &PlanController{service: service}
}

func (p *PlanController) CreatePlan(ctx context.Context, in *planpb.CreatePlanReq) (*planpb.CreatePlanResp, error) {

	err := p.service.CreatePlan(ctx, in.UserId, in.PlanName, &in.ExerciseNames)
	if err != nil {
		return nil, err
	}

	resp := planpb.CreatePlanResp{
		PlanName:      in.PlanName,
		ExerciseNames: in.ExerciseNames,
		Message:       fmt.Sprintf("%v created successfully", in.PlanName),
	}

	return &resp, nil
}

func (p *PlanController) GetAllPlans(ctx context.Context, in *planpb.GetAllPlansReq) (*planpb.GetAllPlansResp, error) {

	numberOfPlans, allPlans, err := p.service.GetAllPlansSer(ctx, int(in.UserId))
	if err != nil {
		return nil, err
	}

	resp := planpb.GetAllPlansResp{}

	for _, eachPlan := range *allPlans {
		eachPlan := planpb.PlanResp{
			PlanName:      eachPlan.PlanName,
			ExerciseNames: eachPlan.Exercises,
		}

		resp.AllPlans = append(resp.AllPlans, &eachPlan)
	}

	resp.NumberOfPlans = int64(numberOfPlans)

	return &resp, nil
}

func (p *PlanController) GetPlanByName(ctx context.Context, in *planpb.GetPlanByNameReq) (*planpb.PlanResp, error) {

	planName, exerciseNames, err := p.service.GetPlanByNameSer(ctx, int(in.UserId), in.PlanName)
	if err != nil {
		return nil, err
	}

	resp := planpb.PlanResp{
		ExerciseNames: *exerciseNames,
		PlanName:      planName,
	}

	return &resp, nil
}

func (p *PlanController) AddExercisesToPlan(ctx context.Context, in *planpb.PlanReq) (*planpb.PlanResp, error) {

	r, err := p.service.AddExercisesToPlan(ctx, int(in.UserId), in.PlanName, &in.ExerciseNames)
	if err != nil {
		return nil, err
	}

	resp := planpb.PlanResp{
		PlanName:      r.PlanName,
		ExerciseNames: r.Exercises,
	}

	return &resp, nil
}

// func (p *PlanController) PlanExistsReturnId(ctx context.Context, in *planpb.SendPlanName) (*planpb.PlanExistsResp, error) {

// 	exits, planId, err := p.service.PlanExistsReturnId(ctx, int(in.UserId), in.PlanName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp := planpb.PlanExistsResp{
// 		Exists: exits,
// 		PlanId: int64(planId),
// 	}

// 	return &resp, nil
// }

// func (p *PlanController) GetEmptyPlanId(ctx context.Context, in *planpb.SendUserID) (*planpb.EmptyPlanIdResp, error) {

// 	emptyPlanId, err := p.service.GetEmptyPlanId(ctx, int(in.UserId))

// 	if err != nil {
// 		return nil, err
// 	}

// 	resp := planpb.EmptyPlanIdResp{
// 		EmptyPlanId: int64(emptyPlanId),
// 	}

// 	return &resp, nil
// }

// func (p *PlanController) PlanExistsReturnPlan(ctx context.Context, in *planpb.SendPlanName) (*planpb.PlanExistsReturnPlanResp, error) {

// 	exists, planId, exerciseIds, err := p.service.PlanExistsReturnPlan(ctx, int(in.UserId), in.PlanName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp := planpb.PlanExistsReturnPlanResp{
// 		Exists: exists,
// 		PlanId: int64(planId),
// 	}

// 	for _, eachExerciseId := range *exerciseIds {
// 		resp.ExerciseIds = append(resp.ExerciseIds, int64(eachExerciseId))
// 	}

// 	return &resp, nil
// }

// func (p *PlanController) CreateEmptyPlan(ctx context.Context, in *planpb.SendUserID) (*planpb.CreateEmptyPlanResp, error) {

// 	err := p.service.CreateEmptyPlan(ctx, int(in.UserId))
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp := planpb.CreateEmptyPlanResp{
// 		Message: "empty plan created successfully",
// 	}

// 	return &resp, nil
// }

// func (p *PlanController) DeleteExercisesFromPlan(ctx context.Context, in *planpb.PlanReq) (*planpb.PlanResp, error) {
// 	r, err := p.service.DeleteExerciseFromPlan(ctx, int(in.UserId), in.PlanName, &in.ExerciseNames)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp := planpb.PlanResp{
// 		PlanName: r.PlanName,
// 		ExerciseNames: r.Exercises,
// 	}

// 	return &resp, nil

// }

// func (p *PlanController) DeletePlan(ctx context.Context, in *planpb.DeletePlanReq) (*planpb.DeletePlanResp, error) {

// 	err := p.service.DeletePlanSer(ctx, int(in.UserId), in.PlanName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &planpb.DeletePlanResp{}, nil
// }

// func (a *PlanController) PING(ctx context.Context, in *planpb.PingPlanReq) (*planpb.PingPlanResp, error) {

// 	return &planpb.PingPlanResp{}, nil
// }

// func (a *PlanController) GetHealth(ctx context.Context, in *planpb.GetHealthReq) (*planpb.GetHealthResp, error) {

// 	pgRespTime, redisRespTime := a.service.GetHealth(ctx)

// 	resp := planpb.GetHealthResp{
// 		PostgresRespTime: durationpb.New(*pgRespTime),
// 		RedisRespTime: durationpb.New(*redisRespTime),
// 	}

// 	return &resp, nil
// }
