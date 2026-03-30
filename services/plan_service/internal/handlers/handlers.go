package handlers

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	customerrors "plan_service/internal/custom_errors"
// 	"plan_service/internal/models"
// 	"plan_service/internal/user"

// 	// "plan_service/internal/middleware"
// 	// "plan_service/internal/models"
// 	"plan_service/internal/services"
// 	"plan_service/internal/utils"

// 	// "plan_service/internal/validators"
// 	"time"

// 	token "wt/pkg/shared"
// )

// // should take the service

// type Handler struct {
// 	service *services.Service
// }

// func NewHandler(s *services.Service) *Handler {
// 	return &Handler{service: s}
// }

// func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {

// 	resp := map[string]string{
// 		"message": "server is alive",
// 	}

// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(resp)
// }

// func (h *Handler) CreatePlan(w http.ResponseWriter, r *http.Request) {

// 	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
// 	defer cancel()

// 	var plan user.Plan2
// 	json.NewDecoder(r.Body).Decode(&plan)

// 	validationErrs, errOccured := plan.Validate()
// 	if errOccured {
// 		utils.ValidationErrWriter(w, validationErrs)
// 		return
// 	}

// 	t := token.JwtToken{}

// 	claims, ok := t.GetClaimsFromContext(ctx)
// 	if !ok {
// 		utils.BadReqErrorWriter(w, customerrors.ErrGettingClaims.Error())
// 		return
// 	}

// 	err := h.service.CreatePlan(ctx, claims.UserId, plan.PlanName, &plan.Exercises)
// 	if err != nil {
// 		if errors.Is(err, customerrors.ErrPlanAlreadyExists) {
// 			customerrors.CustomErrorWriter(w, &customerrors.ErrPlanAlreadyExists2)
// 			return
// 		}
// 		utils.BadReqErrorWriter(w, err.Error())
// 		return
// 	}

// 	resp := models.Plan2Resp{
// 		Message:   fmt.Sprintf("%v created successfully", plan.PlanName),
// 		PlanName:  plan.PlanName,
// 		Exercises: plan.Exercises,
// 	}

// 	utils.CreatedResponseWriter(w, resp)
// }
// func (h *Handler) GetAllPlans(w http.ResponseWriter, r *http.Request) {

// 	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
// 	defer cancel()

// 	t := token.JwtToken{}

// 	claims, ok := t.GetClaimsFromContext(ctx)
// 	if !ok {
// 		utils.BadReqErrorWriter(w, customerrors.ErrGettingClaims.Error())
// 		return
// 	}

// 	resp, err := h.service.GetAllPlansSer(ctx, claims.UserId)
// 	if err != nil {
// 		utils.BadReqErrorWriter(w, err.Error())
// 		return
// 	}

// 	utils.OkResponseWriter(w, resp)
// }

// func (h *Handler) GetPLanByName(w http.ResponseWriter, r *http.Request) {

// 	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
// 	defer cancel()

// 	var planName user.PlanName

// 	json.NewDecoder(r.Body).Decode(&planName)

// 	validationErrs, errOccured := planName.Validate()
// 	if errOccured {
// 		utils.ValidationErrWriter(w, validationErrs)
// 		return
// 	}

// 	t := token.JwtToken{}

// 	claims, ok := t.GetClaimsFromContext(ctx)
// 	if !ok {
// 		utils.BadReqErrorWriter(w, customerrors.ErrGettingClaims.Error())
// 		return
// 	}

// 	resp, err := h.service.GetPlanByNameSer(ctx, claims.UserId, planName.PlanName)
// 	if err != nil {
// 		utils.BadReqErrorWriter(w, err.Error())
// 		return
// 	}

// 	utils.OkResponseWriter(w, resp)
// }

// func (h *Handler) AddExercisesToPlan(w http.ResponseWriter, r *http.Request) {
// 	// req :
// 	// plan_name : "x"
// 	// exericises : [
// 	//    "a", "b", "c", "d", "e"
// 	// ]
// 	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
// 	defer cancel()

// 	var userInput user.Plan2

// 	json.NewDecoder(r.Body).Decode(&userInput)

// 	validationErrs, errOccured := userInput.Validate()
// 	if errOccured {
// 		utils.ValidationErrWriter(w, validationErrs)
// 		return
// 	}

// 	t := token.JwtToken{}

// 	claims, ok := t.GetClaimsFromContext(ctx)
// 	if !ok {
// 		utils.BadReqErrorWriter(w, customerrors.ErrGettingClaims.Error())
// 		return
// 	}

// 	resp, err := h.service.AddExercisesToPlan(ctx, claims.UserId, userInput.PlanName, &userInput.Exercises)
// 	if err != nil {
// 		utils.BadReqErrorWriter(w, err.Error())
// 		return
// 	}

// 	utils.CreatedResponseWriter(w, resp)

// }

// func (h *Handler) DeleteExerciseFromPlan(w http.ResponseWriter, r *http.Request) {
// 	// req :
// 	// plan_name : "x"
// 	// exericises : [
// 	//    "a", "b"
// 	// ]

// 	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
// 	defer cancel()

// 	var userInput user.Plan2

// 	json.NewDecoder(r.Body).Decode(&userInput)

// 	validationErrs, errOccured := userInput.Validate()
// 	if errOccured {
// 		utils.ValidationErrWriter(w, validationErrs)
// 		return
// 	}

// 	t := token.JwtToken{}

// 	claims, ok := t.GetClaimsFromContext(ctx)
// 	if !ok {
// 		utils.BadReqErrorWriter(w, customerrors.ErrGettingClaims.Error())
// 		return
// 	}

// 	resp, err := h.service.DeleteExerciseFromPlan(ctx, claims.UserId, userInput.PlanName, &userInput.Exercises)
// 	if err != nil {
// 		utils.BadReqErrorWriter(w, err.Error())
// 		return
// 	}

// 	utils.CreatedResponseWriter(w, resp)

// }
// func (h *Handler) DeletePlan(w http.ResponseWriter, r *http.Request) {
// 	// req :
// 	// plan_name : "x"
// 	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
// 	defer cancel()

// 	var userInput user.PlanName

// 	json.NewDecoder(r.Body).Decode(&userInput)

// 	validationErrs, errOccured := userInput.Validate()
// 	if errOccured {
// 		utils.ValidationErrWriter(w, validationErrs)
// 		return
// 	}

// 	t := token.JwtToken{}

// 	claims, ok := t.GetClaimsFromContext(ctx)
// 	if !ok {
// 		utils.BadReqErrorWriter(w, customerrors.ErrGettingClaims.Error())
// 		return
// 	}

// 	err := h.service.DeletePlanSer(ctx, claims.UserId, userInput.PlanName)
// 	if err != nil {
// 		utils.BadReqErrorWriter(w, err.Error())
// 		return
// 	}

// 	utils.DeletedRespWriter(w)
// }
