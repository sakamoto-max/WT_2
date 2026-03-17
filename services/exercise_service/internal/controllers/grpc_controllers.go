package controllers

import (
	"context"
	"exercise_service/internal/services"
	exerpb "workout-tracker/proto/shared/exercise"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExerController struct {
	service *services.Service
	exerpb.UnimplementedExerciseServiceServer
}

func NewExerController(service *services.Service) *ExerController {
	return &ExerController{service: service}
}

func (e *ExerController) GetAllExercises(ctx context.Context, in *exerpb.GetAllExercisesREq) (*exerpb.GetAllExercisesResp, error) {
	resp := exerpb.GetAllExercisesResp{}
	r, err := e.service.GetAllExercisesSer(ctx)
	if err != nil{
		return &resp, err
	}

	for _, v := range *r{
		eachExer := exerpb.OneExerciseResp{}
		eachExer.Id = int64(v.Id)
		eachExer.Name = v.Name
		eachExer.BodyPart = v.BodyPart
		eachExer.Equipment = v.Equipment
		eachExer.RestTime = int64(v.RestTime)

		resp.AllExericses = append(resp.AllExericses, &eachExer)
	}

	return &resp, nil
}

func (e *ExerController) GetOneExercise(ctx context.Context, in *exerpb.SendExerciseName) (*exerpb.OneExerciseResp, error) {
	resp := exerpb.OneExerciseResp{}

	r, err := e.service.GetExerciseByNameSer(ctx, in.ExerciseName)
	if err != nil{
		return &resp, err
	}

	resp.Id = int64(r.Id)
	resp.Name = r.Name
	resp.BodyPart = r.BodyPart
	resp.CreatedAt = timestamppb.New(r.CreatedAt)
	resp.RestTime = int64(r.RestTime)
	resp.Equipment = r.Equipment

	return &resp, nil
}


func (e *ExerController) CreateExercise(ctx context.Context, in *exerpb.CreateExerciseReq) (*exerpb.CreateExerciseResp, error) {

	resp := exerpb.CreateExerciseResp{}
	r, err := e.service.CreateExerciseSer(ctx, in.ExerciseName, in.BodyPart, in.Equipment, int(in.RestTime))
	if err != nil{
		return &resp, err
	}

	resp.Message = r.Message
	resp.Exercise.Name = r.Exercise.Name
	// resp.Exercise.Id = int64(r.Exercise.Id)
		resp.Exercise.RestTime = int64(r.Exercise.RestTime)
	resp.Exercise.BodyPart = r.Exercise.BodyPart
	resp.Exercise.Equipment = r.Exercise.Equipment
	resp.Exercise.CreatedAt = timestamppb.New(r.Exercise.CreatedAt)

	return &resp, nil
}
func (e *ExerController) DeleteExercise(ctx context.Context, in *exerpb.SendExerciseName) (*exerpb.DeleteExerciseResp, error) {
	resp := exerpb.DeleteExerciseResp{}
	err := e.service.DeleteExeciseSer(ctx, in.ExerciseName)
	if err != nil{
		return &resp, err
	}

	return &resp, nil
	
}

func (e *ExerController) ExerciseExistsReturnId(ctx context.Context, in *exerpb.SendExerciseName) (*exerpb.ExerciseExistsReturnIdResp, error) {
	var resp exerpb.ExerciseExistsReturnIdResp

	exists, id, err := e.service.ExerciseExistsReturnId(ctx, in.ExerciseName)
	if err != nil{
		return &resp, err
	}

	resp = exerpb.ExerciseExistsReturnIdResp{
		Exists: exists,
		ExerciseId: id,
	}

	return &resp, nil
}
func (e *ExerController) GetExerciseName(ctx context.Context, in *exerpb.SendExerciseID) (*exerpb.GetExerciseNameResp, error) {
		
	var resp exerpb.GetExerciseNameResp
	exerciseName, err := e.service.GetExerciseNameByID(ctx, int(in.ExerciseId))
	if err != nil{
		return &resp, err
	}

	resp.ExerciseName = exerciseName

	return &resp, nil
}

func (a *ExerController) PING(ctx context.Context, in *exerpb.PingExerReq) (*exerpb.PingExerResp, error) {
	r := exerpb.PingExerResp{}

	return &r, nil
}