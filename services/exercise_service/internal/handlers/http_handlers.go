package handlers

import (
	"context"
	"encoding/json"
	"exercise_service/internal/models"
	"exercise_service/internal/services"
	"exercise_service/internal/utils"
	"fmt"
	"net/http"
	"time"
)

type Handler struct {
	service *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{
		service: s,
	}
}



func (h *Handler) GetAllExercises(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	resp, err := h.service.GetAllExercisesSer(ctx)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, resp, http.StatusOK)
}

func (h *Handler) GetExerciseByName(w http.ResponseWriter, r *http.Request) {
	// pathvalue to json
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	exerciseName := r.PathValue("exerciseName")

	resp, err := h.service.GetExerciseByNameSer(ctx, exerciseName)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, resp, http.StatusOK)
}

// needs validation
func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()
	
	var exercise models.Exercise

	json.NewDecoder(r.Body).Decode(&exercise)

	resp, err := h.service.CreateExerciseSer(ctx, &exercise)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, resp, http.StatusCreated)
}

// complete the handler
// needs validation
func (h *Handler) UpdateExercise(w http.ResponseWriter, r *http.Request) {

}

// pathvalue to json
// needs validation
func (h *Handler) DeleteExecise(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	exerciseName := r.PathValue("exerciseName")

	// exerciseName := r.PathValue("exerciseName")

	err := h.service.DeleteExeciseSer(ctx, exerciseName)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	resp := models.DelExerciseResp{}

	resp.Message = fmt.Sprintf("exercise deleted successfully : %v", exerciseName)

	utils.ResponseWriter(w, resp, http.StatusOK)
}
