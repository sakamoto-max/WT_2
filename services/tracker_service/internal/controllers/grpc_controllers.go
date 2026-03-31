package controllers

import (
	"context"
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

	tracker := models.Tracker{}

	for _, eachExerWithidSetsandReps := range in.AllExerices {
		a := models.Workout{}
		a.ExerciseId = int(eachExerWithidSetsandReps.ExerciseId)
		for _, eachExer := range a.RepsWeight {
			reps := models.Reps{}
			reps.Reps = int(eachExer.Reps)
			reps.Weight = int(eachExer.Weight)

			a.RepsWeight = append(a.RepsWeight, reps)
		}

		tracker.Workout = append(tracker.Workout, a)
	}

	err := t.service.EndWorkoutSer(ctx, in.UserId, &tracker)
	if err != nil {
		return nil, err
	}

	resp := trackerpb.EndWorkoutResp{
		Message: "workout ended successfully",
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
