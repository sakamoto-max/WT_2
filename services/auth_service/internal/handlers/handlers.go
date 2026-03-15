package handlers

import (
	// customerrors "auth_service/internal/custom_errors"
	"auth_service/internal/responses"
	"auth_service/internal/services"
	"auth_service/internal/user"
	"auth_service/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	myErrs "wt/pkg/my_errors"

	// myerrors "wt/pkg/my_errors"
	pkg "wt/pkg/shared"
	myutils "wt/pkg/utils"
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
func (h *Handler) UserSignUp(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.Signup

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	createdAt, err := h.service.SignUp(ctx, userInput.Name, userInput.Email, userInput.Password)
	if err != nil {
		switch {
		case errors.Is(err, myErrs.ErrNameAlreadyExits):
			myutils.BadReq(w, err)
		case errors.Is(err, myErrs.ErrEmailAlreadyExits):
			myutils.BadReq(w, err)
		default:
			log.Printf("server encountered an err : %v", err)
			myutils.InternalServerErr(w, err)
		}
		return
	}

	resp := responses.SignUpResp{
		Name:      userInput.Name,
		Email:     userInput.Email,
		Role:      "user",
		CreatedAt: createdAt,
	}

	utils.CreatedRespWriter(w, resp)
}

// needs validation
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {

	var userInput user.Login

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	resp, err := h.service.Login(ctx, userInput.Email, userInput.Password)

	if err != nil {
		switch {
		case errors.Is(err, myErrs.ErrEmailNotFound):
			myutils.BadReq(w, err)
		case errors.Is(err, myErrs.ErrIncorrectPassword):
			myutils.BadReq(w, err)
		default:
			myutils.InternalServerErr(w, err)
		}
		// utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, resp, http.StatusOK)
}
func (h *Handler) UserLogOut(w http.ResponseWriter, r *http.Request) {

	t := pkg.JwtToken{}
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

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var uuid user.UUIDReader
	json.NewDecoder(r.Body).Decode(&uuid)

	validationErrs, errOccured := uuid.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
	}

	accessToken, err := h.service.GetNewAccessTokenSer(ctx, uuid.UUID)
	if err != nil {
		switch {
		case errors.Is(err, myErrs.ErrTokenExpired):
			err := myErrs.NewAppErr(myErrs.ErrTokenExpired, http.StatusBadRequest)
			err.AppErrWriter(w)
		default:
			utils.InternalServerErr(w, err)
		}
		return
	}

	resp := map[string]string{
		"new_access_token": accessToken,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
