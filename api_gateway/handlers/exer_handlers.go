package handlers

import (
	"api_gateway/user"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	exerpb "workout-tracker/proto/shared/exercise"
	myerrors "wt/pkg/my_errors"
	token "wt/pkg/jwt"
	"wt/pkg/utils"
)

func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

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
		UserId:       claims.UserId,
	}

	res, err := h.exerClient.CreateExercise(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	resp := user.CreateExerResp{
		Messsage: fmt.Sprintf("exercise %v has been successfully created ", in.ExerciseName),
		Exercise: user.Exercise{
			Id:        res.Id,
			Name:      in.ExerciseName,
			BodyPart:  in.BodyPart,
			Equipment: in.Equipment,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	utils.CreatedWriter(w, resp)
}
func (h *Handler) GetExerciseByName(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.ExerciseName

	json.NewDecoder(r.Body).Decode(&userInput)
	fmt.Println(userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	fmt.Println("making the in")

	in := exerpb.SendExerciseName{
		ExerciseName: userInput.Name,
		UserId:       claims.UserId,
	}

	fmt.Println("sending the in")
	res, err := h.exerClient.GetOneExercise(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	resp := user.Exercise{
		Id:        res.Id,
		Name:      res.Name,
		BodyPart:  res.BodyPart,
		Equipment: res.Equipment,
		CreatedAt: res.CreatedAt.AsTime(),
		UpdatedAt: res.UpdatedAt.AsTime(),
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) GetAllExercises(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	in := exerpb.GetAllExercisesREq{UserId: claims.UserId}

	res, err := h.exerClient.GetAllExercises(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
	}

	var resp []user.Exercise

	for _, eachExer := range res.AllExericses {
		exer := user.Exercise{
			Id:        eachExer.Id,
			Name:      eachExer.Name,
			BodyPart:  eachExer.BodyPart,
			Equipment: eachExer.Equipment,
			CreatedAt: eachExer.CreatedAt.AsTime(),
			UpdatedAt: eachExer.UpdatedAt.AsTime(),
		}

		resp = append(resp, exer)
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) DeleteExecise(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.ExerciseName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	in := exerpb.SendExerciseName{
		ExerciseName: userInput.Name,
		UserId: claims.UserId,
	}

	resp, err := h.exerClient.DeleteExercise(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	utils.DeletedNotFoundWriter(w, resp)
}
func (h *Handler) UpdateExercise(w http.ResponseWriter, r *http.Request) {}