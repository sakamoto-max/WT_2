package services

import (
	"auth_service/internal/jwt"
	"auth_service/internal/mappings"
	"auth_service/internal/utils"
	"context"
	"strings"

	"github.com/google/uuid"
	myerrs "github.com/sakamoto-max/wt_2_pkg/myerrs"
	pb "github.com/sakamoto-max/wt_2_proto/shared/auth"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) UserSignUp(ctx context.Context, in *pb.UserSignUpReq) (*pb.UserSignUpResp, error) {

	in.Email = strings.ToLower(in.Email)

	hashedPass, err := utils.HashThePassword(in.Password)
	if err != nil {
		return nil, err
	}

	in.Password = hashedPass

	userId, createdAt, err := s.db.Auth.CreateUser(ctx, mappings.ToUserSignUp(in))
	if err != nil {
		return nil, err
	}

	return &pb.UserSignUpResp{
		Name:      in.Name,
		Email:     in.Email,
		Role:      in.Role,
		CreatedAt: timestamppb.New(createdAt),
		UserId:    userId,
	}, nil

}

func (s *Service) UserLogin(ctx context.Context, in *pb.UserLoginReq) (*pb.UserLoginResp, error) {

	hashedPass, err := s.db.Password.FetchUserPass(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	err = utils.MatchPasswords(in.Password, hashedPass)
	if err != nil {
		return nil, err
	}

	userId, roleID, name, err := s.db.UserDeatails.FetchUserIdRoleIdName(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	AccessToken, err := jwt.GenerateAccessToken(userId, roleID)
	if err != nil {
		return nil, err
	}

	exists, err := s.cache.Token.RefreshExists(ctx, userId)
	if err != nil {
		return nil, err
	}

	if exists {
		UUID, err := s.cache.Uuid.GetUUID(ctx, userId)
		if err != nil {
			return nil, err
		}

		return &pb.UserLoginResp{
			Name:        name,
			UserId:      userId,
			Message:     "login successful",
			AccessToken: AccessToken,
			UUID:        UUID,
			Email:       in.Email,
		}, nil
	}

	refreshToken, err := jwt.GenerateRefreshToken(userId, roleID)
	if err != nil {
		return nil, err
	}
	UUID := uuid.NewString()
	err = s.cache.Token.SetRefreshTokenAndUUID(ctx, UUID, refreshToken, userId)
	if err != nil {
		return nil, err
	}

	return &pb.UserLoginResp{
		Name:        name,
		UserId:      userId,
		Message:     "login successful",
		AccessToken: AccessToken,
		UUID:        UUID,
		Email:       in.Email,
	}, nil
}

func (s *Service) UserLogOut(ctx context.Context, in *pb.SendUserId) (*pb.UserLogOutResp, error) {
	uuid, err := s.cache.Uuid.GetUUID(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	err = s.cache.User.UserLogout(ctx, in.UserId, uuid)
	if err != nil {
		return nil, err
	}

	return &pb.UserLogOutResp{
		Message: "logout successful",
	}, nil
}

func (s *Service) GetNewAccessToken(ctx context.Context, in *pb.SendUUID) (*pb.GetNewAccessTokenResp, error) {
	refreshToken, err := s.cache.Token.GetRefreshToken(ctx, in.UUID)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := jwt.GenerateAccessToken(claims.UserId, claims.RoleId)
	if err != nil {
		return nil, err
	}

	return &pb.GetNewAccessTokenResp{
		NewAccessToken: accessToken,
	}, nil
}

func (s *Service) PING(ctx context.Context, in *pb.PINGreq) (*pb.PINGresp, error) {

	return &pb.PINGresp{}, nil
}

func (s *Service) ChangePass(ctx context.Context, in *pb.ChangePassReq) (*pb.ChangePassResp, error) {
	if in.OldPass == in.NewPass {
		return nil, myerrs.BadReqErrMaker(ErrSameOldPass)
	}

	passFromDb, err := s.db.Password.FetchUserPassById(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	err = utils.MatchPasswords(in.OldPass, passFromDb)
	if err != nil {
		return nil, myerrs.BadReqErrMaker(ErrPasswordIncorrect)
	}

	hashedNewPass, err := utils.HashThePassword(in.NewPass)
	if err != nil {
		return nil, err
	}

	err = s.db.Password.ChangePass(ctx, in.UserId, hashedNewPass)
	if err != nil {
		return nil, err
	}

	return &pb.ChangePassResp{
		Message: "password changed successfully",
	}, nil
}

func (s *Service) ChangeEmail(ctx context.Context, in *pb.ChangeEmailReq) (*pb.ChangeEmailResp, error) {
	if in.OldEmail == in.NewEmail {
		return nil, myerrs.BadReqErrMaker(ErrSameOldEmail)
	}

	emailFromDb, err := s.db.Email.GetEmail(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	if emailFromDb != in.OldEmail {
		return nil, myerrs.BadReqErrMaker(ErrEmailDoesntMatch)
	}

	err = s.db.Email.ChangeEmail(ctx, in.UserId, in.NewEmail)
	if err != nil {
		return nil, err
	}

	return &pb.ChangeEmailResp{
		Message: "email changed successfully",
	}, nil
}
func (s *Service) GetHealth(ctx context.Context, in *pb.GetHealthReq) (*pb.GetHealthResp, error) {
	pgRespTime := s.db.Metrics.GetRespTime(ctx)
	redisRespTime := s.cache.Metrics.GetRespTime(ctx)

	if pgRespTime == nil && redisRespTime == nil {
		return &pb.GetHealthResp{
			PostgresRespTime: nil,
			RedisRespTime:    nil,
		}, nil
	}
	if redisRespTime == nil {
		return &pb.GetHealthResp{
			PostgresRespTime: durationpb.New(*pgRespTime),
			RedisRespTime:    nil,
		}, nil
	}
	if pgRespTime == nil {
		return &pb.GetHealthResp{
			PostgresRespTime: nil,
			RedisRespTime:    durationpb.New(*redisRespTime),
		}, nil
	}

	return &pb.GetHealthResp{
		PostgresRespTime: durationpb.New(*pgRespTime),
		RedisRespTime:    durationpb.New(*redisRespTime),
	}, nil
}
