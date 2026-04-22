package handlers

import (
	"api_gateway/responses"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	trackpb "workout-tracker/proto/shared/tracker"
	"wt/pkg/user"

	// token "wt/pkg/jwt"
	"wt/pkg/middleware"
	myerrors "wt/pkg/my_errors"
	"wt/pkg/utils"

	"go.uber.org/zap"
)

func (h *Handler) StartEmptyWorkout(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

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

	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

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

	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("END_WORKOUT_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.Tracker

	json.NewDecoder(r.Body).Decode(&userInput)

	if userInput.UserResponse != "" {

		in := trackpb.EndWorkoutReq{}

		switch userInput.UserResponse {
		case "yes":
			in.UserResponse = true
		case "no":
			in.UserResponse = false
		}

		in.UserId = claims.UserId

		resp, err := h.trackClient.EndWorkout(ctx, &in)
		if err != nil {
			myerrors.ErrMatcher(w, err)
			return
		}

		if resp.ConflictOccured {
			resp := responses.TrackerConflict{
				RequestStatus: resp.RequestStatus,
				Reason:        resp.Reason,
				ExerciseNames: resp.ExerciseNames,
				Message:       resp.Message,
			}

			utils.ConflictWriter(w, resp)
			logger.Log.Infow("END_WORKOUT_CONFLICT_OCCURED", zap.String("REQ_ID", reqId))
			return
		}

		utils.OkRespWriter(w, resp)

		logger.Log.Infow("END_WORKOUT_TRACKED", zap.String("REQ_ID", reqId))
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
			rw.Weight = eachExer.Weight

			allExer.SetsAndReps = append(allExer.SetsAndReps, &rw)
		}

		in.AllExerices = append(in.AllExerices, &allExer)
	}

	switch userInput.UserResponse {
	case "yes":
		in.UserResponse = true
	case "no":
		in.UserResponse = false
	}

	resp, err := h.trackClient.EndWorkout(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		// utils.BadReq(w, err)
		return
	}

	if resp.ConflictOccured {
		resp := responses.TrackerConflict{
			RequestStatus: resp.RequestStatus,
			Reason:        resp.Reason,
			ExerciseNames: resp.ExerciseNames,
			Message:       resp.Message,
		}

		utils.ConflictWriter(w, resp)
		logger.Log.Infow("END_WORKOUT_CONFLICT_OCCURED", zap.String("REQ_ID", reqId))
		return
	}

	utils.OkRespWriter(w, resp)

	logger.Log.Infow("END_WORKOUT_TRACKED", zap.String("REQ_ID", reqId))
	// validationErrs, errOccured := userInput.Validate()
	// if errOccured {
	// 	utils.ValidationErrWriter(w, *validationErrs)
	// 	return
	// }
}

func (h *Handler) CancelWorkout(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("CANCEL_WORKOUT_CALLED", zap.String("REQ_ID", reqId))

	in := trackpb.CancelWorkoutReq{
		UserId: claims.UserId,
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	resp, err := h.trackClient.CancelWorkout(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	utils.OkRespWriter(w, resp)

	logger.Log.Infow("CANCEL_WORKOUT_SUCCESSFULL", zap.String("REQ_ID", reqId))
}

// {
// 	"exercises" : [
// 		{
// 			"exercise_name" : "x",
// 			"tracker" : [
// 				{
// 					"weight" : 20,
// 					"reps" : 10
// 				},
// 				{
// 					"weight" : 20,
// 					"reps" : 10
// 				},
// 				{
// 					"weight" : 20,
// 					"reps" : 10
// 				}
// 			]
// 		},
// 		{
// 			"exercise_name" : "y",
// 			"tracker" : [
// 				{
// 					"weight" : 20,
// 					"reps" : 10
// 				},
// 				{
// 					"weight" : 20,
// 					"reps" : 10
// 				},
// 				{
// 					"weight" : 20,
// 					"reps" : 10
// 				}
// 			]
// 		},
// 		{
// 			"exercise_name" : "z",
// 			"tracker" : [
// 				{
// 					"weight" : 20,
// 					"reps" : 10
// 				},
// 				{
// 					"weight" : 20,
// 					"reps" : 10
// 				},
// 				{
// 					"weight" : 20,
// 					"reps" : 10
// 				}
// 			]
// 		}
// 	]
// }
