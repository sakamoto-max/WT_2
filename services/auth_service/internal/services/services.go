package services

import (
	customerrors "auth_service/internal/custom_errors"
	"auth_service/internal/repository"
	"auth_service/internal/responses"
	"auth_service/internal/utils"
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type action string

var (
	userSignedUp action = "user_signed_up"
)

type Service struct {
	repo      *repository.Repo
	// mqChannel *amqp.Channel
}

func NewService(r *repository.Repo) *Service {
	return &Service{repo: r}
}

func (s *Service) SignUp(ctx context.Context, name string, email string, password string) (*responses.SignUpResp, error) {

	var resp *responses.SignUpResp

	email = strings.ToLower(email)
	exists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		return resp, err
	}

	if exists {
		return resp, customerrors.ErrUserAlreadyExists
	}

	hashedPass, err := utils.HashThePassword(password)
	if err != nil {
		return resp, err
	}

	resp, err = s.repo.SignUp(ctx, name, email, hashedPass)
	if err != nil {
		return resp, err
	}

	resp.Role = "user"

	// err = Produce(ctx, string(userSignedUp), userId, "plan", s.mqChannel)
	// if err != nil {
	// 	return resp, fmt.Errorf("error in producing to the queue : %w/n", err)
	// }

	return resp, nil
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

	hashedPass, err := s.repo.FetchHashedPass(ctx, email)
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

	token := utils.JwtToken{}
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

	err = s.repo.Logout(ctx, userId, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetNewAccessTokenSer(ctx context.Context, UUID uuid.UUID) (string, error) {

	refreshToken, err := s.repo.GetRefreshToken(ctx, UUID)
	if err != nil {
		return "", fmt.Errorf("error getting token from repo : %w", err)
	}

	token := utils.JwtToken{}

	claims, err := token.ValidateToken(refreshToken)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			return "", customerrors.ErrRefreshTokenExp
		}
		//
		return "", fmt.Errorf("error validating the token : %w", err)
	}

	accessToken, err := token.GenerateAccessToken(claims.UserId, claims.RoleId)
	if err != nil {
		return "", err
	}

	return accessToken, nil
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
