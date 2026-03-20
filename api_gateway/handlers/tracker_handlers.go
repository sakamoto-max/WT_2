package handlers

import (
	"api_gateway/user"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	trackpb "workout-tracker/proto/shared/tracker"
	myerrors "wt/pkg/my_errors"
	token "wt/pkg/shared"
	"wt/pkg/utils"
)

func (h *Handler) StartEmptyWorkout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	fmt.Printf("user Id : %v\n", claims.UserId)

	in := trackpb.StartEmptyWorkoutReq{
		UserId: int64(claims.UserId),
	}

	resp, err := h.trackClient.StartEmptyWorkout(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	utils.CreatedWriter(w, resp)
}
func (h *Handler) StartWorkoutWithPlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.PlanName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := trackpb.StartWorkoutWithPlanReq{
		UserId:   int64(claims.UserId),
		PlanName: userInput.PlanName,
	}

	resp, err := h.trackClient.StartWorkoutWithPlan(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.CreatedWriter(w, resp)
}
func (h *Handler) EndWorkout(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.Tracker

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := trackpb.EndWorkoutReq{}
	in.UserId = int64(claims.UserId)

	for _, allExercises := range userInput.Workout {
		allExer := trackpb.TrackerForEachExer{}
		allExer.ExerciseId = int64(allExercises.ExerciseId)

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
}
