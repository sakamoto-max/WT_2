package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	planpb "workout-tracker/proto/shared/plan"
	"wt/pkg/middleware"
	myerrors "wt/pkg/my_errors"
	"wt/pkg/user"
	"wt/pkg/utils"

	"go.uber.org/zap"
)

func (h *Handler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	logger := middleware.GetLogger(r.Context())
	reqId := middleware.GetReqId(r.Context())

	logger.Log.Infow("CREATE_PLAN_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
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
	claims := middleware.GetClaims(r.Context())
	logger := middleware.GetLogger(r.Context())
	reqId := middleware.GetReqId(r.Context())

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

	claims := middleware.GetClaims(r.Context())
	logger := middleware.GetLogger(r.Context())
	reqId := middleware.GetReqId(r.Context())

	logger.Log.Infow("GET_PLAN_BY_NAME_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.PlanName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
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

	claims := middleware.GetClaims(r.Context())
	logger := middleware.GetLogger(r.Context())
	reqId := middleware.GetReqId(r.Context())

	logger.Log.Infow("ADD_EXERCISE_TO_PLAN_CALLED", zap.String("REQ_ID", reqId))
	
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	
	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)
	
	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
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
	claims := middleware.GetClaims(r.Context())
	logger := middleware.GetLogger(r.Context())
	reqId := middleware.GetReqId(r.Context())
	
	
	logger.Log.Infow("DELETE_EXERCISE_CALLED", zap.String("REQ_ID", reqId))
	
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	

	var userInput user.Plan
	
	json.NewDecoder(r.Body).Decode(&userInput)
	
	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
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

	claims := middleware.GetClaims(r.Context())
	logger := middleware.GetLogger(r.Context())
	reqId := middleware.GetReqId(r.Context())

	logger.Log.Infow("DELETE_EXERCISE_SUCCESSFULL", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()



	var userInput user.PlanName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
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
