package controllers

import (
	"auth_service/internal/services"
	"context"

	pb "workout-tracker/proto/shared/auth"
	myerrors "wt/pkg/my_errors"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthController struct {
	pb.UnimplementedAuthServiceServer
	service *services.Service
}

func NewAuthController(service *services.Service) *AuthController {
	return &AuthController{service: service}
}

func (a *AuthController) UserSignUp(ctx context.Context, in *pb.UserSignUpReq) (*pb.UserSignUpResp, error) {

	r := pb.UserSignUpResp{}

	createdAt, err := a.service.SignUp(ctx, in.Name, in.Email, in.Password, in.Role)
	if err != nil {
		err = myerrors.ErrMaker(err)
		return &r, err
	}

	r.Name = in.Name
	r.Email = in.Email
	r.Role = in.Role
	r.CreatedAt = timestamppb.New(createdAt)

	return &r, nil
}

func (a *AuthController) UserLogin(ctx context.Context, in *pb.UserLoginReq) (*pb.UserLoginResp, error) {

	r := pb.UserLoginResp{}

	resp, err := a.service.Login(ctx, in.Email, in.Password)
	if err != nil {
		err = myerrors.ErrMaker(err)
		return &r, err
	}

	r.Message = "login successful"
	r.Email = in.Email
	r.Name = resp.Name
	r.AccessToken = resp.AccessToken
	r.UUID = resp.UUID

	return &r, nil
}

// errs possible :
// 1. user not logged in
// 2. interna server error
func (a *AuthController) UserLogOut(ctx context.Context, in *pb.SendUserId) (*pb.UserLogOutResp, error) {
	r := pb.UserLogOutResp{}
	err := a.service.Logout(ctx, int(in.UserId))
	if err != nil {
		return &r, err
	}

	r.Message = "logout successful"
	return &r, nil
}

// possible errs :
// 1. not logged in
// 2.
func (a *AuthController) GetNewAccessToken(ctx context.Context, in *pb.SendUUID) (*pb.GetNewAccessTokenResp, error) {
	r := pb.GetNewAccessTokenResp{}

	accessToken, err := a.service.GetNewAccessTokenSer(ctx, in.UUID)
	if err != nil {
		return &r, err
	}

	r.NewAccessToken = accessToken
	return &r, nil
}

func (a *AuthController) PING(ctx context.Context, in *pb.PINGreq) (*pb.PINGresp, error) {
	r := pb.PINGresp{}

	return &r, nil
}

func (a *AuthController) ChangePass(ctx context.Context, in *pb.ChangePassReq) (*pb.ChangePassResp, error) {

	r := pb.ChangePassResp{}

	err := a.service.ChangePass(ctx, int(in.UserId), in.OldPass, in.NewPass)
	if err != nil {
		return &r, err
	}

	r.Message = "password changed successfully"

	return &r, nil

}

func (a *AuthController) ChangeEmail(ctx context.Context, in *pb.ChangeEmailReq) (*pb.ChangeEmailResp, error) {

	r := pb.ChangeEmailResp{}
	err := a.service.ChangeEmail(ctx, int(in.UserId), in.OldEmail, in.NewEmail)
	if err != nil {
		err = myerrors.ErrMaker(err)
		return &r, err
	}

	r.Message = "email changed successfully"

	return &r, nil
}

func (a *AuthController) GetHealth(ctx context.Context, in *pb.GetHealthReq) (*pb.GetHealthResp, error) {

	resp := pb.GetHealthResp{}

	pgRespTime, redisRespTime := a.service.GetHealth(ctx)

	resp.PostgresRespTime = durationpb.New(*pgRespTime)
	resp.RedisRespTime = durationpb.New(*redisRespTime)

	return &resp, nil
}