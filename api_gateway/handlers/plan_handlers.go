package handlers

import (
	"api_gateway/user"
	"context"
	"encoding/json"
	"net/http"
	"time"
	myerrors "wt/pkg/my_errors"
	token "wt/pkg/jwt"
	"wt/pkg/utils"
	planpb "workout-tracker/proto/shared/plan"
)

func (h *Handler) CheckHealthPlan(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := planpb.CreatePlanReq{
		UserId:        claims.UserId,
		PlanName:      userInput.PlanName,
		ExerciseNames: userInput.Exercises,
	}

	resp, err := h.planClient.CreatePlan(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	utils.CreatedWriter(w, resp)
}


func (h *Handler) GetAllPlans(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	in := planpb.GetAllPlansReq{
		UserId: claims.UserId,
	}

	resp, err := h.planClient.GetAllPlans(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) GetPLanByName(w http.ResponseWriter, r *http.Request) {

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

	in := planpb.GetPlanByNameReq{
		UserId:   claims.UserId,
		PlanName: userInput.PlanName,
	}

	resp, err := h.planClient.GetPlanByName(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) AddExercisesToPlan(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := planpb.PlanReq{
		UserId:        claims.UserId,
		PlanName:      userInput.PlanName,
		ExerciseNames: userInput.Exercises,
	}

	resp, err := h.planClient.AddExercisesToPlan(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) DeleteExerciseFromPlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := planpb.PlanReq{
		UserId:        claims.UserId,
		PlanName:      userInput.PlanName,
		ExerciseNames: userInput.Exercises,
	}

	resp, err := h.planClient.DeleteExercisesFromPlan(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) DeletePlan(w http.ResponseWriter, r *http.Request) {
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

	in := planpb.DeletePlanReq{
		UserId:   claims.UserId,
		PlanName: userInput.PlanName,
	}

	resp, err := h.planClient.DeletePlan(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
