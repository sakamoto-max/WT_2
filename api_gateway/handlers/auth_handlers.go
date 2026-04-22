package handlers

import (
	"api_gateway/responses"
	"context"
	"encoding/json"
	"net/http"
	"time"
	authpb "workout-tracker/proto/shared/auth"
	"wt/pkg/enum"
	"wt/pkg/middleware"
	myerrors "wt/pkg/my_errors"
	"wt/pkg/user"
	"wt/pkg/utils"

	"go.uber.org/zap"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	logger.Log.Infow("USER_SIGNUP_CALLED", zap.String("REQ_ID", reqId))

	var userInput user.Signup

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
		return
	}

	in := authpb.UserSignUpReq{}
	in.Email = userInput.Email
	in.Name = userInput.Name
	in.Password = userInput.Password

	if userInput.Role == nil {
		in.Role = string(enum.UserRole)
	} else {
		in.Role = *userInput.Role
	}

	resp, err := h.authClient.UserSignUp(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher2(w, err, logger)
		return
	}

	re := responses.SignUpResp{
		Name:      resp.Name,
		Email:     resp.Email,
		Role:      resp.Role,
		CreatedAt: resp.CreatedAt.AsTime(),
	}

	utils.CreatedWriter(w, re)

	logger.Log.Infow("USER_SIGNUP_SUCCESSFUL", zap.String("REQ_ID", reqId))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("USER_LOGIN_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.Login

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErr)
		return
	}

	in := authpb.UserLoginReq{
		Email:    userInput.Email,
		Password: userInput.Password,
	}

	resp, err := h.authClient.UserLogin(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
	logger.Log.Infow("USER_LOGIN_SUCCESFULL", zap.String("REQ_ID", reqId))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("USER_LOGOUT_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	in := authpb.SendUserId{
		UserId: claims.UserId,
	}

	resp, err := h.authClient.UserLogOut(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
	logger.Log.Infow("USER_LOGOUT_SUCCESSFULL", zap.String("REQ_ID", reqId))

}

func (h *Handler) GetNewAccessToken(w http.ResponseWriter, r *http.Request) {

	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("NEW_ACCESS_TOKEN_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.UUIDReader

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErr)
		return
	}

	in := authpb.SendUUID{
		UUID: userInput.UUID,
	}

	resp, err := h.authClient.GetNewAccessToken(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.CreatedWriter(w, resp)

	logger.Log.Infow("NEW_ACCESS_TOKEN_CREATED", zap.String("REQ_ID", reqId))
}

func (h *Handler) ChangePassWord(w http.ResponseWriter, r *http.Request) {

	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("USER_PASSWORD_CHANGE_CALLED", zap.String("REQ_ID", reqId))
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.ChangePass

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
		return
	}

	in := authpb.ChangePassReq{
		UserId:  claims.UserId,
		OldPass: userInput.OldPass,
		NewPass: userInput.NewPass,
	}

	resp, err := h.authClient.ChangePass(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
	logger.Log.Infow("USER_PASSWORD_CHANGE_SUCCESSFULL", zap.String("REQ_ID", reqId))
}

func (h *Handler) ChangeEmail(w http.ResponseWriter, r *http.Request) {

	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("USER_EMAIL_CHANGE_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.ChangeEmail

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErrs)
		return
	}

	in := authpb.ChangeEmailReq{
		UserId:   claims.UserId,
		OldEmail: userInput.OldEmail,
		NewEmail: userInput.NewEmail,
	}
	resp, err := h.authClient.ChangeEmail(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
	logger.Log.Infow("USER_EMAIL_CHANGE_SUCCESSFULL", zap.String("REQ_ID", reqId))
}
