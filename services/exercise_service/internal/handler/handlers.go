package handler

import (
	"context"
	"exercise_service/internal/services"
	exerpb "github.com/sakamoto-max/wt_2-proto/shared/exercise" 
	"github.com/sakamoto-max/wt_2-pkg/logger"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	service *services.Service
	exerpb.UnimplementedExerciseServiceServer
	logger *logger.MyLogger
}

func NewHandler(service *services.Service, logger *logger.MyLogger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (e *Handler) GetAllExercises(ctx context.Context, in *exerpb.GetAllExercisesREq) (*exerpb.GetAllExercisesResp, error) {
	resp := exerpb.GetAllExercisesResp{}
	allExers, err := e.service.GetAllExercisesSer(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	for _, exer := range *allExers {
		eachExer := exerpb.OneExerciseResp{
			Id:        exer.Id,
			Name:      exer.Name,
			BodyPart:  exer.BodyPart,
			Equipment: exer.Equipment,
			CreatedAt: timestamppb.New(exer.CreatedAt),
			UpdatedAt: timestamppb.New(exer.UpdatedAt),
		}

		resp.AllExericses = append(resp.AllExericses, &eachExer)
	}

	resp.NumberOfExercises = int64(len(*allExers))

	return &resp, nil
}
func (e *Handler) GetOneExercise(ctx context.Context, in *exerpb.SendExerciseName) (*exerpb.OneExerciseResp, error) {

	r, err := e.service.GetExerciseByName(ctx, in.UserId, in.ExerciseName)
	if err != nil {
		return nil, err
	}
	
	resp := exerpb.OneExerciseResp{
		Id:        r.Id,
		Name:      r.Name,
		BodyPart:  r.BodyPart,
		Equipment: r.Equipment,
		CreatedAt: timestamppb.New(r.CreatedAt),
		UpdatedAt: timestamppb.New(r.UpdatedAt),
	}

	return &resp, nil
}
func (e *Handler) CreateExercise(ctx context.Context, in *exerpb.CreateExerciseReq) (*exerpb.CreateExerciseResp, error) {

	UUID, err := e.service.CreateExerciseSer(ctx, in.UserId, in.ExerciseName, in.BodyPart, in.Equipment)
	if err != nil {
		return nil, err
	}

	resp := exerpb.CreateExerciseResp{Id: UUID}
	return &resp, nil
}
func (e *Handler) DeleteExercise(ctx context.Context, in *exerpb.SendExerciseName) (*exerpb.DeleteExerciseResp, error) {
	resp := exerpb.DeleteExerciseResp{}
	err := e.service.DeleteExeciseSer(ctx, in.UserId, in.ExerciseName)
	if err != nil {
		return nil, err
	}

	return &resp, nil

}
func (e *Handler) ExerciseExistsReturnId(ctx context.Context, in *exerpb.SendExerciseName) (*exerpb.ExerciseExistsReturnIdResp, error) {
	var resp exerpb.ExerciseExistsReturnIdResp

	id, err := e.service.ExerciseExistsReturnId(ctx, in.UserId, in.ExerciseName)
	if err != nil {
		return &resp, err
	}

	resp = exerpb.ExerciseExistsReturnIdResp{
		ExerciseId: id,
	}

	return &resp, nil
}
func (e *Handler) GetExerciseName(ctx context.Context, in *exerpb.SendExerciseID) (*exerpb.GetExerciseNameResp, error) {

	var resp exerpb.GetExerciseNameResp
	exerciseName, err := e.service.GetExerciseNameByID(ctx, in.ExerciseId)
	if err != nil {
		return &resp, err
	}

	resp.ExerciseName = exerciseName

	return &resp, nil
}
func (a *Handler) PING(ctx context.Context, in *exerpb.PingExerReq) (*exerpb.PingExerResp, error) {
	r := exerpb.PingExerResp{}

	return &r, nil
}
func (a *Handler) GetHealth(ctx context.Context, in *exerpb.GetHealthReq) (*exerpb.GetHealthResp, error) {

	resp := exerpb.GetHealthResp{}

	pgRespTime, redisRespTime := a.service.GetHealth(ctx)

	resp.PostgresRespTime = durationpb.New(*pgRespTime)
	resp.RedisRespTime = durationpb.New(*redisRespTime)

	return &resp, nil
}
