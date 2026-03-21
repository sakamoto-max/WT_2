package controllers

import (
	"context"
	"fmt"
	"tracker_service/internal/services"
	"tracker_service/internal/user"
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

	resp := trackerpb.GeneralResp{}

	err := t.service.StartEmptyWorkoutSer(ctx, int(in.UserId))
	if err != nil {
		if err == myerrors.ErrWorkoutOngoing {
			st := status.New(codes.AlreadyExists, err.Error())
			return &resp, st.Err()
		}
		return &resp, err
	}

	resp.Message = "an empty workout has started"

	return &resp, nil
}
func (t *TrackerController) StartWorkoutWithPlan(ctx context.Context, in *trackerpb.StartWorkoutWithPlanReq) (*trackerpb.StartWorkoutWithPlanResp, error) {
	resp := trackerpb.StartWorkoutWithPlanResp{}

	exercisesNames, err := t.service.StartWorkoutWithPlanSer(ctx, int(in.UserId), in.PlanName)
	if err != nil {
		return &resp, err
	}

	resp.Message = fmt.Sprintf("workout with plan %v has started", in.PlanName)
	resp.PlanName = in.PlanName
	resp.ExercisesInPlan = *exercisesNames

	return &resp, nil
}
func (t *TrackerController) EndWorkout(ctx context.Context, in *trackerpb.EndWorkoutReq) (*trackerpb.EndWorkoutResp, error) {
	resp := trackerpb.EndWorkoutResp{}

	tracker := user.Tracker{}
	// workout := user.Workout{}
	// reps := user.Reps{}

	for _, eachExerWithidSetsandReps := range in.AllExerices {
		a := user.Workout{}
		a.ExerciseId = int(eachExerWithidSetsandReps.ExerciseId)
		for _, eachExer := range a.RepsWeight {
			reps := user.Reps{}
			reps.Reps = int(eachExer.Reps)
			reps.Weight = int(eachExer.Weight)

			a.RepsWeight = append(a.RepsWeight, reps)
		}

		tracker.Workout = append(tracker.Workout, a)
	}

	err := t.service.EndWorkoutSer(ctx, int(in.UserId), &tracker)
	if err != nil {
		return &resp, err
	}

	resp.Message = "workout ended successfully"

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
