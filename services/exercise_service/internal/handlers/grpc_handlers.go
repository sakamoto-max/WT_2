package handlers

import (
	"context"
	"exercise_service/internal/services"
)

import exercisepb "workout-tracker/proto/shared/exercise"

type ExerciseHandler struct {
	exercisepb.UnimplementedExerciseServiceServer
	service *services.ExerciseService
}

func NewExerciseHandler(service *services.ExerciseService) *ExerciseHandler {
	
	return &ExerciseHandler{
		service: service,
	}
}

func (e *ExerciseHandler) ExerciseExistsReturnId(ctx context.Context, req *exercisepb.SendExerciseName) (*exercisepb.ExerciseExistsReturnIdResp, error) {

	var resp exercisepb.ExerciseExistsReturnIdResp

	exists, id, err := e.service.ExerciseExistsReturnId(ctx, req.ExerciseName)
	if err != nil{
		return &resp, err
	}

	resp = exercisepb.ExerciseExistsReturnIdResp{
		Exists: exists,
		ExerciseId: id,
	}

	return &resp, nil
}