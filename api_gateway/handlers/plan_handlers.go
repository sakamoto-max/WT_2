package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	planpb "github.com/sakamoto-max/wt_2-proto/shared/plan"
	"github.com/sakamoto-max/wt_2-pkg/middleware"
	myerrors "github.com/sakamoto-max/wt_2-pkg/my_errors"
	"github.com/sakamoto-max/wt_2/api_gateway/user"
	"go.uber.org/zap"
	"github.com/sakamoto-max/wt_2/api_gateway/utils"
)

func (h *Handler) CreatePlan(w http.ResponseWriter, r *http.Request) {
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

	logger.Log.Infow("CREATE_PLAN_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	err = userInput.Validate()
	if err != nil {
		user.ValidationErrWriter(w, err)
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

	logger.Log.Infow("CREATE_PLAN_SUCCESSFULL", zap.String("REQ_ID", reqId))
}
func (h *Handler) GetAllPlans(w http.ResponseWriter, r *http.Request) {
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

	logger.Log.Infow("GET_ALL_PLANS_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	in := planpb.GetAllPlansReq{
		UserId: claims.UserId,
	}

	resp, err := h.planClient.GetAllPlans(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
	logger.Log.Infow("GET_ALL_PLANS_SUCCESSFULL", zap.String("REQ_ID", reqId))
}
func (h *Handler) GetPLanByName(w http.ResponseWriter, r *http.Request) {

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

	logger.Log.Infow("GET_PLAN_BY_NAME_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.PlanName

	json.NewDecoder(r.Body).Decode(&userInput)

	err = userInput.Validate()
	if err != nil {
		user.ValidationErrWriter(w, err)
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
	logger.Log.Infow("GET_PLAN_BY_NAME_SUCCESSFULL", zap.String("REQ_ID", reqId))
}
func (h *Handler) AddExercisesToPlan(w http.ResponseWriter, r *http.Request) {

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

	logger.Log.Infow("ADD_EXERCISE_TO_PLAN_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	err = userInput.Validate()
	if err != nil {
		user.ValidationErrWriter(w, err)
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
	logger.Log.Infow("ADD_EXERCISE_TO_PLAN_SUCCESSFULL", zap.String("REQ_ID", reqId))
}
func (h *Handler) DeleteExerciseFromPlan(w http.ResponseWriter, r *http.Request) {
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

	logger.Log.Infow("DELETE_EXERCISE_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	err = userInput.Validate()
	if err != nil {
		user.ValidationErrWriter(w, err)
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
	logger.Log.Infow("DELETE_EXERCISE_SUCCESSFULL", zap.String("REQ_ID", reqId))
}
func (h *Handler) DeletePlan(w http.ResponseWriter, r *http.Request) {

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

	logger.Log.Infow("DELETE_EXERCISE_SUCCESSFULL", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.PlanName

	json.NewDecoder(r.Body).Decode(&userInput)

	err = userInput.Validate()
	if err != nil {
		user.ValidationErrWriter(w, err)
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
