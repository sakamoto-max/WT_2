package controllers

import (
	"context"
	"errors"
	"fmt"
	"tracker_service/internal/models"
	"tracker_service/internal/services"
	trackerpb "workout-tracker/proto/shared/tracker"
	myerrors "wt/pkg/my_errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type TrackerController struct {
	service *services.Service
	trackerpb.UnimplementedTrackerServiceServer
}

func NewTrackerController(service *services.Service) *TrackerController {
	return &TrackerController{
		service: service,
	}
}

func (t *TrackerController) StartEmptyWorkout(ctx context.Context, in *trackerpb.StartEmptyWorkoutReq) (*trackerpb.GeneralResp, error) {

	err := t.service.StartEmptyWorkoutSer(ctx, in.UserId)
	if err != nil {
		if err == myerrors.ErrWorkoutOngoing {
			st := status.New(codes.AlreadyExists, err.Error())
			return nil, st.Err()
		}
		return nil, err
	}

	resp := trackerpb.GeneralResp{
		Message: "an empty workout has started",
	}

	return &resp, nil
}
func (t *TrackerController) StartWorkoutWithPlan(ctx context.Context, in *trackerpb.StartWorkoutWithPlanReq) (*trackerpb.StartWorkoutWithPlanResp, error) {

	exercisesNames, err := t.service.StartWorkoutWithPlanSer(ctx, in.UserId, in.PlanName)
	if err != nil {
		if err == myerrors.ErrWorkoutOngoing {
			st := status.New(codes.AlreadyExists, err.Error())
			return nil, st.Err()
		}
		return nil, err
	}

	resp := trackerpb.StartWorkoutWithPlanResp{
		Message:         fmt.Sprintf("workout with plan %v has started", in.PlanName),
		PlanName:        in.PlanName,
		ExercisesInPlan: *exercisesNames,
	}
	return &resp, nil
}
func (t *TrackerController) EndWorkout(ctx context.Context, in *trackerpb.EndWorkoutReq) (*trackerpb.EndWorkoutResp, error) {

	tracker := convertToLocal(in)

	msg, err := t.service.EndWorkoutSer(ctx, in.UserId, &tracker)
	if err != nil {
		var target *myerrors.Conflict
		if errors.As(err, &target){

			resp := trackerpb.EndWorkoutResp{
				RequestStatus: target.RequestStatus,
				Message: target.Message,
				Reason: target.Reason.Error(),
				ExerciseNames: target.ExerciseNames,
				ConflictOccured: true,
			}

			return &resp, nil
		}

		return nil, err
	}
	
	resp := trackerpb.EndWorkoutResp{}

	switch msg {
	case nil:
		resp.Message = "workout ended successfully"
	default:
		resp.Message = *msg
	}

	return &resp, nil
}

func (a *TrackerController) PING(ctx context.Context, in *trackerpb.PingTrackReq) (*trackerpb.PingTrackResp, error) {
	r := trackerpb.PingTrackResp{}

	return &r, nil
}


func (a *TrackerController) GetHealth(ctx context.Context, in *trackerpb.GetHealthReq) (*trackerpb.GetHealthResp, error) {

	resp := trackerpb.GetHealthResp{}

	pgRespTime, redisRespTime := a.service.GetHealth(ctx)

	resp.PostgresRespTime = durationpb.New(*pgRespTime)
	resp.RedisRespTime = durationpb.New(*redisRespTime)

	return &resp, nil
}

func (a *TrackerController) CancelWorkout(ctx context.Context, in *trackerpb.CancelWorkoutReq) (*trackerpb.CancelWorkoutResp, error) {

	err := a.service.CancelWorkout(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	resp := trackerpb.CancelWorkoutResp{
		Message: "workout has been successfully canceled",
	}

	return &resp, nil
}



func convertToLocal(in *trackerpb.EndWorkoutReq) models.Tracker {

	main := models.Tracker{}

	if in.UserResponse {
		main.UserResponse = in.UserResponse
		return main
	}

	for _, eachExer := range in.AllExerices {
		w := models.Workout{
			ExerciseName: eachExer.ExerciseName,
		}

		for _, repsPlusWeight := range eachExer.SetsAndReps {

			rPlusW := models.Reps{
				Reps:   int(repsPlusWeight.Reps),
				Weight: int(repsPlusWeight.Weight),
			}

			w.RepsWeight = append(w.RepsWeight, rPlusW)
		}
		main.Workout = append(main.Workout, w)
	}

	return main
}
