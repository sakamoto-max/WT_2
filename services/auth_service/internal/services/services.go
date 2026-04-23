package services

import (
	"auth_service/internal/repository"
	"auth_service/internal/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	// token "wt/pkg/jwt"
	token "github.com/sakamoto-max/wt_2-pkg/jwt"
	myerrors "github.com/sakamoto-max/wt_2-pkg/my_errors"
	// myerrors "wt/pkg/my_errors"
	"github.com/google/uuid"
)

type Service struct {
	repo *repository.Repo
}

func NewService(r *repository.Repo) *Service {
	return &Service{repo: r}
}

var (
	ErrEmailDoesntMatch  = errors.New("email user sent is wrong")
	ErrPasswordIncorrect = errors.New("incorrect password")
	ErrSameOldPass       = errors.New("new password cannot be same as the old password")
	ErrSameOldEmail      = errors.New("new email cannot be same as the old email")
)

func (s *Service) SignUp(ctx context.Context, name string, email string, password string, role string) (string, time.Time, error) {

	var CreatedAt time.Time

	email = strings.ToLower(email)

	hashedPass, err := utils.HashThePassword(password)
	if err != nil {
		return "", CreatedAt, err
	}

	userId, CreatedAt, err := s.repo.CreateUser(ctx, name, email, hashedPass, role)
	if err != nil {
		return "", CreatedAt, err
	}

	return userId, CreatedAt, nil
}
func (s *Service) Login(ctx context.Context, email string, password string) (string, string, string, string, error) {

	var refreshToken string
	var UUID string

	hashedPass, err := s.repo.FetchUserPass(ctx, email)
	if err != nil {
		return "", "", "", "", err
	}

	err = utils.MatchPasswords(password, hashedPass)
	if err != nil {
		return "", "", "", "", err
	}

	userId, roleID, name, err := s.repo.FetchUserIdRoleIdName(ctx, email)
	if err != nil {
		return "", "", "", "", err
	}

	token := token.JwtToken{}
	AccessToken, err := token.GenerateAccessToken(userId, roleID)
	if err != nil {
		return "", "", "", "", err
	}

	exists, err := s.repo.RefreshExists(ctx, userId)
	if err != nil {
		return "", "", "", "", err
	}

	if exists {
		UUID, err = s.repo.GetUUID(ctx, userId)
		if err != nil {
			return "", "", "", "", err
		}
	}

	if !exists {
		refreshToken, err = token.GenerateRefreshToken(userId, roleID)
		if err != nil {
			return "", "", "", "", err
		}
		UUID = uuid.NewString()
		err := s.repo.SetRefreshTokenAndUUID(ctx, UUID, refreshToken, userId)
		if err != nil {
			return "", "", "", "", err
		}
	}
	return userId, name, AccessToken, UUID, nil
}
func (s *Service) Logout(ctx context.Context, userId string) error {

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
		return "", err
	}

	accessToken, err := token.GenerateAccessToken(claims.UserId, claims.RoleId)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
func (s *Service) ChangePass(ctx context.Context, userId string, oldPass string, newPass string) error {

	if oldPass == newPass {
		return myerrors.BadReqErrMaker(ErrSameOldPass)
	}

	passFromDb, err := s.repo.FetchUserPassById(ctx, userId)
	if err != nil {
		return err
	}

	err = utils.MatchPasswords(oldPass, passFromDb)
	if err != nil {
		return myerrors.BadReqErrMaker(ErrPasswordIncorrect)
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
}
func (s *Service) ChangeEmail(ctx context.Context, userId string, oldEmail string, newEmail string) error {

	if oldEmail == newEmail {
		return myerrors.BadReqErrMaker(ErrSameOldEmail)
	}

	emailFromDb, err := s.repo.GetEmail(ctx, userId)
	if err != nil {
		return err
	}

	if emailFromDb != oldEmail {
		return myerrors.BadReqErrMaker(ErrEmailDoesntMatch)
	}

	err = s.repo.ChangeEmail(ctx, userId, newEmail)
	if err != nil {
		return err
	}

	return nil
}
func (s *Service) GetHealth(ctx context.Context) (*time.Duration, *time.Duration) {

	pgRespTime := s.repo.GetPostgresRespTime(ctx)
	redisRespTime := s.repo.GetRedisRespTime(ctx)

	return pgRespTime, redisRespTime
}