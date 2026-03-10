package handlers

import (
	"auth_service/internal/middleware"
	"auth_service/internal/models"
	"auth_service/internal/services"
	"auth_service/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// the user should only get error code and
// {
// 	"message" : "internal server error"
// }
// errors specifics should be logged in the terminal

type Handler struct {
	service *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{service: s}
}

// needs validation
// name, email, password
func (h *Handler) UserSignUp(w http.ResponseWriter, r *http.Request) {

	userDetails := r.Context().Value(middleware.UserSignUp).(*models.Signup)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	resp, err := h.service.SignUp(ctx, userDetails.Name, userDetails.Email, userDetails.Password)
	if err != nil {
		fmt.Printf("error occured : %v\n", err)
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, resp, http.StatusCreated)
}

// needs validation
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {

	userDetails := r.Context().Value(middleware.UserLogin).(*models.Login)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	resp, err := h.service.Login(ctx, userDetails.Email, userDetails.Password)

	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, resp, http.StatusOK)
}

func (h *Handler) UserLogOut(w http.ResponseWriter, r *http.Request) {

	t := utils.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())

	if !ok {
		response := map[string]string{
			"message": "error getting claims from context",
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	err := h.service.Logout(ctx, claims.UserId)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"message": "logout successful",
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

// needs validation
func (h *Handler) GetNewAccessToken(w http.ResponseWriter, r *http.Request) {

	uuid := r.Context().Value(middleware.RefreshTokenUUID).(*models.UUIDReader)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	accessToken, err := h.service.GetNewAccessTokenSer(ctx, uuid.UUID)
	if err != nil {
		fmt.Printf("error making the token : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"new_access_token": accessToken,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
