package services

import (
	"auth_service/internal/repository"
	"context"
	"time"
)

type service struct {
	repo repository.RepoIface
}

func NewService(r repository.RepoIface) ServiceIface {
	return &service{repo: r}
}

type ServiceIface interface {
	SignUp(ctx context.Context, name string, email string, password string, role string) (string, time.Time, error)
	Login(ctx context.Context, email string, password string) (string, string, string, string, error)
	Logout(ctx context.Context, userId string) error
	GetNewAccessTokenSer(ctx context.Context, UUID string) (string, error)
	ChangePass(ctx context.Context, userId string, oldPass string, newPass string) error
	ChangeEmail(ctx context.Context, userId string, oldEmail string, newEmail string) error
	GetHealth(ctx context.Context) (*time.Duration, *time.Duration)
}
