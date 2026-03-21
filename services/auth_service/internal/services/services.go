package services

import (
	customerrors "auth_service/internal/custom_errors"
	"auth_service/internal/repository"
	"auth_service/internal/responses"
	"auth_service/internal/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	planpb "workout-tracker/proto/shared/plan"
	myerrors "wt/pkg/my_errors"
	token "wt/pkg/shared"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// type action string

// var (
//
//	userSignedUp action = "user_signed_up"
//
// )
type Service struct {
	repo       *repository.Repo
	planClient planpb.PlanServiceClient
	// mqChannel *amqp.Channel
}

func NewService(r *repository.Repo, planClient planpb.PlanServiceClient) *Service {
	return &Service{repo: r, planClient: planClient}
}

func (s *Service) SignUp(ctx context.Context, name string, email string, password string, role string) (time.Time, error) {

	var CreatedAt time.Time

	email = strings.ToLower(email)

	hashedPass, err := utils.HashThePassword(password)
	if err != nil {
		return CreatedAt, err
	}

	_, CreatedAt, err = s.repo.CreateUser(ctx, name, email, hashedPass, role)
	if err != nil {
		return CreatedAt, err
	}

	// _, err = s.planClient.CreateEmptyPlan(ctx, &planpb.SendUserID{UserId: int64(userId),})
	// if err != nil{
	// 	// st,  := status.FromError(err)

	// 	return CreatedAt, myerrors.PlanServerNotResponding

	// 	// try again
	// 	// return CreatedAt, fmt.Errorf("error creating empty plan : %v", err)
	// 	// grpc.ErrClientConnClosing
	// 	// return CreatedAt, status.Newf(codes.Canceled, "plan server is not responding").Err()
	// }

	return CreatedAt, nil
}

func (s *Service) Login(ctx context.Context, email string, password string) (*responses.LoginResp, error) {

	var refreshToken string
	var UUID string

	var resp responses.LoginResp

	exists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		return &resp, err
	}

	if !exists {
		return &resp, customerrors.PleaseSignUp
	}

	hashedPass, err := s.repo.FetchUserPass(ctx, email)
	if err != nil {
		return &resp, err
	}

	err = utils.MatchPasswords(password, hashedPass)
	if err != nil {
		return &resp, err
	}

	userId, roleID, name, err := s.repo.FetchNameUserIdRoleId(ctx, email)
	if err != nil {
		return &resp, err
	}

	token := token.JwtToken{}
	AccessToken, err := token.GenerateAccessToken(userId, roleID)
	if err != nil {
		return &resp, err
	}

	exists, err = s.repo.RefreshExists(ctx, userId)
	if err != nil {
		return &resp, err
	}

	if exists {
		UUID, err = s.repo.GetUUID(ctx, userId)
		if err != nil {
			return &resp, err
		}
	}

	if !exists {
		refreshToken, err = token.GenerateRefreshToken(userId, roleID)
		if err != nil {
			return &resp, err
		}
		UUID = uuid.NewString()
		err := s.repo.SetRefreshTokenAndUUID(ctx, UUID, refreshToken, userId)
		if err != nil {
			return &resp, err
		}
	}

	resp.Message = "login successful"
	resp.Email = email
	resp.Name = name
	resp.AccessToken = AccessToken
	resp.UUID = UUID
	return &resp, nil
}

func (s *Service) Logout(ctx context.Context, userId int) error {

	// remove the refresh token from redis

	uuid, err := s.repo.GetUUID(ctx, userId)
	if err != nil {
		return err
	}

	err = s.repo.UserLogout(ctx, userId, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetNewAccessTokenSer(ctx context.Context, UUID string) (string, error) {

	refreshToken, err := s.repo.GetRefreshToken(ctx, UUID)
	if err != nil {
		return "", fmt.Errorf("error getting token from repo : %w", err)
	}

	token := token.JwtToken{}

	claims, err := token.ValidateToken(refreshToken)
	if err != nil {
		if errors.Is(err, myerrors.ErrTokenExpired) {
			return "", myerrors.ErrRefreshExpired
		}
		return "", err
	}

	accessToken, err := token.GenerateAccessToken(claims.UserId, claims.RoleId)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *Service) ChangePass(ctx context.Context, userId int, oldPass string, newPass string) error {
	// {
	// 	"old_password" : "x",
	// 	"new_password" : "y",
	// }

	if oldPass == newPass {
		return myerrors.ErrOldPassNewPassSame
	}

	passFromDb, err := s.repo.FetchUserPassById(ctx, userId)
	if err != nil {
		return err
	}

	err = utils.MatchPasswords(oldPass, passFromDb)
	if err != nil {
		return myerrors.ErrIncorrectPassword
	}

	hashedNewPass, err := utils.HashThePassword(newPass)
	if err != nil {
		return err
	}

	err = s.repo.ChangePass(ctx, userId, hashedNewPass)
	if err != nil {
		return err
	}

	return nil

	// check if ui old_pass and new_pass r same -> old_pass cannot be same as new pass
	// get the old pass from the db
	// check if the ui old_pass and the pass from the db r same
	// if not -> incorrect_old pass
	// if yes -> successfully changed the password
}

func (s *Service) ChangeEmail(ctx context.Context, userId int, oldEmail string, newEmail string) error {

	if oldEmail == newEmail {
		return myerrors.ErrOldEmailNewEmailSame
	}

	emailFromDb, err := s.repo.GetEmail(ctx, userId)
	if err != nil {
		return err
	}

	if emailFromDb != oldEmail {
		return myerrors.ErrEmailDoesntMatch
	}

	err = s.repo.ChangeEmail(ctx, userId, newEmail)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetHealth(ctx context.Context) (*time.Duration, *time.Duration) {

	// check resp time of pg

	pgRespTime := s.repo.GetPostgresRespTime(ctx)
	redisRespTime := s.repo.GetRedisRespTime(ctx)

	return pgRespTime, redisRespTime
}

// func Produce(ctx context.Context, action string, userID int, queueName string, ch *amqp.Channel) error {
// 	a := models.MqMsg{
// 		UserId: userID,
// 		Action: action,
// 		Time:   time.Now(),
// 	}

// 	data, err := utils.MakeJSON(a)
// 	if err != nil {
// 		return err
// 	}

// 	msg := amqp.Publishing{
// 		Body: data,
// 	}

// 	// ch.PublishWithContext()
// 	// ch.PublishWithContext(ctx, "", "repo", false, false, msg)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }
