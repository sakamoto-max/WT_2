package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	customerrors "plan_service/internal/custom_errors"
	"plan_service/internal/middleware"
	"plan_service/internal/models"
	"plan_service/internal/services"
	"plan_service/internal/utils"
	"plan_service/internal/validators"
	"time"
)

// should take the service

type Handler struct {
	service *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{service: s}
}


func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {

	resp := map[string]string{
		"message": "server is alive",
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
func (h *Handler) CreatePlan(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var plan models.Plan2
	json.NewDecoder(r.Body).Decode(&plan)

	validationErrs, err := validators.CreatePlanValidator(&plan)
	if err != nil {
		validators.ValidationErrWriter(w, validationErrs)
		return
	}

	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		utils.BadReqErrorWriter(w, customerrors.ErrGettingClaims.Error())
		return
	}

	resp, err := h.service.CreatePlanSer(ctx, claims.UserId, &plan)
	if err != nil {
		if errors.Is(err, customerrors.ErrPlanAlreadyExists) {
			customerrors.CustomErrorWriter(w, &customerrors.ErrPlanAlreadyExists2)
			return
		}
		utils.BadReqErrorWriter(w, err.Error())
		return
	}

	utils.CreatedResponseWriter(w, resp)
}
func (h *Handler) GetAllPlans(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		utils.BadReqErrorWriter(w, customerrors.ErrGettingClaims.Error())
		return
	}

	resp, err := h.service.GetAllPlansSer(ctx, claims.UserId)
	if err != nil {
		utils.BadReqErrorWriter(w, err.Error())
		return
	}

	utils.OkResponseWriter(w, resp)
}
func (h *Handler) GetPLanByName(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	planName := r.PathValue("planName")

	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		utils.BadReqErrorWriter(w, customerrors.ErrGettingClaims.Error())
		return
	}

	resp, err := h.service.GetPlanByNameSer(ctx, claims.UserId, planName)
	if err != nil {
		utils.BadReqErrorWriter(w, err.Error())
		return
	}

	utils.OkResponseWriter(w, resp)
}
func (h *Handler) UpdateThePlan(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) DeleteAPlan(w http.ResponseWriter, r *http.Request) {

}
