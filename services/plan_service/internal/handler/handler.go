package handler

import (
	"context"
	"fmt"
	"plan_service/internal/services"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Handler struct {
	service services.ServiceIface
	planpb.UnimplementedPlanServiceServer
	logger *logger.MyLogger
}

func NewHandler(service services.ServiceIface, logger *logger.MyLogger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (p *Handler) CreatePlan(ctx context.Context, in *planpb.CreatePlanReq) (*planpb.CreatePlanResp, error) {

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
func (p *Handler) GetAllPlans(ctx context.Context, in *planpb.GetAllPlansReq) (*planpb.GetAllPlansResp, error) {

	numberOfPlans, allPlans, err := p.service.GetAllPlansSer(ctx, in.UserId)
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
func (p *Handler) GetPlanByName(ctx context.Context, in *planpb.GetPlanByNameReq) (*planpb.GetPlanByNameResp, error) {

	planId, planName, exerciseNames, err := p.service.GetPlanByNameSer(ctx, in.UserId, in.PlanName)
	if err != nil {
		return nil, err
	}

	resp := planpb.GetPlanByNameResp{
		ExerciseNames: *exerciseNames,
		PlanName:      planName,
		PlanId:        planId,
	}

	return &resp, nil
}
func (p *Handler) AddExercisesToPlan(ctx context.Context, in *planpb.PlanReq) (*planpb.PlanResp, error) {

	r, err := p.service.AddExercisesToPlan(ctx, in.UserId, in.PlanName, &in.ExerciseNames)
	if err != nil {
		return nil, err
	}

	resp := planpb.PlanResp{
		PlanName:      r.PlanName,
		ExerciseNames: r.Exercises,
	}

	return &resp, nil
}
func (a *Handler) GetHealth(ctx context.Context, in *planpb.GetHealthReq) (*planpb.GetHealthResp, error) {

	pgRespTime, redisRespTime := a.service.GetHealth(ctx)

	resp := planpb.GetHealthResp{
		PostgresRespTime: durationpb.New(*pgRespTime),
		RedisRespTime:    durationpb.New(*redisRespTime),
	}

	return &resp, nil
}
func (p *Handler) GetEmptyPlanId(ctx context.Context, in *planpb.SendUserID) (*planpb.EmptyPlanIdResp, error) {

	emptyPlanId, err := p.service.GetEmptyPlanId(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	resp := planpb.EmptyPlanIdResp{
		EmptyPlanId: emptyPlanId,
	}

	return &resp, nil
}
func (p *Handler) DeletePlan(ctx context.Context, in *planpb.DeletePlanReq) (*planpb.DeletePlanResp, error) {

	err := p.service.DeletePlanSer(ctx, in.UserId, in.PlanName)
	if err != nil {
		return nil, err
	}

	return &planpb.DeletePlanResp{}, nil
}
func (p *Handler) DeleteExercisesFromPlan(ctx context.Context, in *planpb.PlanReq) (*planpb.PlanResp, error) {
	r, err := p.service.DeleteExerciseFromPlan(ctx, in.UserId, in.PlanName, &in.ExerciseNames)
	if err != nil {
		return nil, err
	}

	resp := planpb.PlanResp{
		PlanName: r.PlanName,
		ExerciseNames: r.Exercises,
	}

	return &resp, nil

}






