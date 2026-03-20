package handlers

import (
	"api_gateway/user"
	"context"
	"encoding/json"
	"net/http"
	"time"
	"wt/pkg/utils"
	exerpb "workout-tracker/proto/shared/exercise"
)

func (h *Handler) GetAllExercises(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	in := exerpb.GetAllExercisesREq{}

	resp, err := h.exerClient.GetAllExercises(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) GetExerciseByName(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.ExerciseName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	in := exerpb.SendExerciseName{
		ExerciseName: userInput.Name,
	}

	resp, err := h.exerClient.GetOneExercise(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.Exercise

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	in := exerpb.CreateExerciseReq{
		ExerciseName: userInput.Name,
		BodyPart:     userInput.BodyPart,
		Equipment:    userInput.Equipment,
		RestTime:     int64(userInput.RestTime),
	}

	resp, err := h.exerClient.CreateExercise(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.CreatedWriter(w, resp)
}
func (h *Handler) DeleteExecise(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.ExerciseName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	in := exerpb.SendExerciseName{
		ExerciseName: userInput.Name,
	}

	resp, err := h.exerClient.DeleteExercise(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.DeletedNotFoundWriter(w, resp)
}
