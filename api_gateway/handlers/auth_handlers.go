package handlers

import (
	"api_gateway/responses"
	"api_gateway/user"
	"context"
	"encoding/json"
	"net/http"
	"time"
	authpb "workout-tracker/proto/shared/auth"
	myerrors "wt/pkg/my_errors"
	token "wt/pkg/jwt"
	"wt/pkg/utils"
	"wt/pkg/enum"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	// parse the data from the user
	var userInput user.Signup

	json.NewDecoder(r.Body).Decode(&userInput)

	// validate it
	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	// make a new req
	in := authpb.UserSignUpReq{}
	in.Email = userInput.Email
	in.Name = userInput.Name
	in.Password = userInput.Password
	
	if userInput.Role == nil{
		in.Role = string(enum.UserRole)
	}else{
		in.Role = *userInput.Role
	}


	resp, err := h.authClient.UserSignUp(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	re := responses.SignUpResp{
		Name:      resp.Name,
		Email:     resp.Email,
		Role:      resp.Role,
		CreatedAt: resp.CreatedAt.AsTime(),
	}

	utils.CreatedWriter(w, re)
}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.Login

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
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
}
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	in := authpb.SendUserId{
		UserId: claims.UserId,
	}

	resp, err := h.authClient.UserLogOut(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) GetNewAccessToken(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.UUIDReader

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
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
}

func (h *Handler) ChangePassWord(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.ChangePass

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
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

	// {
	// 	"old_password" : "x",
	// 	"new_password" : "y",
	// }

	// check if ui old_pass and new_pass r same -> old_pass cannot be same as new pass
	// get the old pass from the db
	// check if the ui old_pass and the pass from the db r same
	// if not -> incorrect_old pass
	// if yes -> successfully changed the password
}
func (h *Handler) ChangeEmail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.ChangeEmail

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
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
}
