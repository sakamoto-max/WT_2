package controllers

import (
	"auth_service/internal/services"
	"context"

	pb "workout-tracker/proto/shared/auth"

	"wt/pkg/logger"
	myerrors "wt/pkg/my_errors"

	// "wt/pkg/middleware"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthController struct {
	pb.UnimplementedAuthServiceServer
	service *services.Service
	logger  *logger.MyLogger
}

func NewAuthController(service *services.Service, logger *logger.MyLogger) *AuthController {
	return &AuthController{service: service, logger: logger}
}

func (a *AuthController) UserSignUp(ctx context.Context, in *pb.UserSignUpReq) (*pb.UserSignUpResp, error) {

	userId, createdAt, err := a.service.SignUp(ctx, in.Name, in.Email, in.Password, in.Role)
	if err != nil {
		return nil, err
	}

	r := pb.UserSignUpResp{
		Name:      in.Name,
		Email:     in.Email,
		Role:      in.Role,
		CreatedAt: timestamppb.New(createdAt),
		UserId:    userId,
	}

	return &r, nil
}

func (a *AuthController) UserLogin(ctx context.Context, in *pb.UserLoginReq) (*pb.UserLoginResp, error) {

	userId, name, accesToken, UUID, err := a.service.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	r := pb.UserLoginResp{
		Message:     "login successful",
		Email:       in.Email,
		Name:        name,
		AccessToken: accesToken,
		UUID:        UUID,
		UserId:      userId,
	}

	return &r, nil
}

func (a *AuthController) UserLogOut(ctx context.Context, in *pb.SendUserId) (*pb.UserLogOutResp, error) {
	err := a.service.Logout(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	r := pb.UserLogOutResp{Message: "logout successful"}

	return &r, nil
}

func (a *AuthController) GetNewAccessToken(ctx context.Context, in *pb.SendUUID) (*pb.GetNewAccessTokenResp, error) {

	accessToken, err := a.service.GetNewAccessTokenSer(ctx, in.UUID)
	if err != nil {
		return nil, err
	}

	r := pb.GetNewAccessTokenResp{NewAccessToken: accessToken}

	return &r, nil
}

func (a *AuthController) PING(ctx context.Context, in *pb.PINGreq) (*pb.PINGresp, error) {
	r := pb.PINGresp{}

	return &r, nil
}

func (a *AuthController) ChangePass(ctx context.Context, in *pb.ChangePassReq) (*pb.ChangePassResp, error) {

	err := a.service.ChangePass(ctx, in.UserId, in.OldPass, in.NewPass)
	if err != nil {
		return nil, err
	}

	r := pb.ChangePassResp{Message: "password changed successfully"}

	return &r, nil

}

func (a *AuthController) ChangeEmail(ctx context.Context, in *pb.ChangeEmailReq) (*pb.ChangeEmailResp, error) {

	err := a.service.ChangeEmail(ctx, in.UserId, in.OldEmail, in.NewEmail)
	if err != nil {
		err = myerrors.ErrMaker(err)
		return nil, err
	}

	r := pb.ChangeEmailResp{Message: "email changed successfully"}

	return &r, nil
}

func (a *AuthController) GetHealth(ctx context.Context, in *pb.GetHealthReq) (*pb.GetHealthResp, error) {

	pgRespTime, redisRespTime := a.service.GetHealth(ctx)

	resp := pb.GetHealthResp{
		PostgresRespTime: durationpb.New(*pgRespTime),
		RedisRespTime:    durationpb.New(*redisRespTime),
	}

	return &resp, nil
}
