package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
		trackpb "workout-tracker/proto/shared/tracker"
	"time"
	"wt/pkg/user"

	// token "wt/pkg/jwt"
	"wt/pkg/middleware"
	myerrors "wt/pkg/my_errors"
	"wt/pkg/utils"

	"go.uber.org/zap"
)

func (h *Handler) StartEmptyWorkout(w http.ResponseWriter, r *http.Request) {

	claims := middleware.GetClaims(r.Context())
	logger := middleware.GetLogger(r.Context())
	reqId := middleware.GetReqId(r.Context())

	logger.Log.Infow("START_EMPTY_WORKOUT_CALLED", zap.String("REQ_ID", reqId))
	
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	fmt.Printf("user Id : %v\n", claims.UserId)
	
	in := trackpb.StartEmptyWorkoutReq{
		UserId: claims.UserId,
	}

	resp, err := h.trackClient.StartEmptyWorkout(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}
	
	utils.CreatedWriter(w, resp)
	logger.Log.Infow("START_EMPTY_WORKOUT_SUCCESSFULL", zap.String("REQ_ID", reqId))
}

func (h *Handler) StartWorkoutWithPlan(w http.ResponseWriter, r *http.Request) {
	
	claims := middleware.GetClaims(r.Context())
	logger := middleware.GetLogger(r.Context())
	reqId := middleware.GetReqId(r.Context())
	
	logger.Log.Infow("START_WORKOUT_WITH_PLAN_CALLED", zap.String("REQ_ID", reqId))
	
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	
	var userInput user.PlanName
	
	json.NewDecoder(r.Body).Decode(&userInput)
	
	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
		return
	}

	in := trackpb.StartWorkoutWithPlanReq{
		UserId:   claims.UserId,
		PlanName: userInput.PlanName,
	}

	resp, err := h.trackClient.StartWorkoutWithPlan(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	utils.CreatedWriter(w, resp)

	logger.Log.Infow("START_WORKOUT_WITH_PLAN_SUCCESSFUL", zap.String("REQ_ID", reqId))
}

func (h *Handler) EndWorkout(w http.ResponseWriter, r *http.Request) {
	
	claims := middleware.GetClaims(r.Context())
	logger := middleware.GetLogger(r.Context())
	reqId := middleware.GetReqId(r.Context())
	
	logger.Log.Infow("END_WORKOUT_CALLED", zap.String("REQ_ID", reqId))
	
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.Tracker

	json.NewDecoder(r.Body).Decode(&userInput)
	
	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
		return
	}

	in := trackpb.EndWorkoutReq{}
	in.UserId = claims.UserId
	
	for _, allExercises := range userInput.Workout {
		allExer := trackpb.TrackerForEachExer{}
		allExer.ExerciseName = allExercises.ExerciseName
		
		for _, eachExer := range allExercises.RepsWeight {
			rw := trackpb.SetsAndReps{}
			rw.Reps = int64(eachExer.Reps)
			rw.Weight = int64(eachExer.Weight)
			
			allExer.SetsAndReps = append(allExer.SetsAndReps, &rw)
		}
		
		in.AllExerices = append(in.AllExerices, &allExer)
	}
	
	resp, err := h.trackClient.EndWorkout(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}
	
	utils.OkRespWriter(w, resp)

	logger.Log.Infow("END_WORKOUT_TRACKED", zap.String("REQ_ID", reqId))
}
