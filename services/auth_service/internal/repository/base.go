package repository

import (
	"auth_service/internal/mappings"
	"context"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbIface interface {
	CreateUser(ctx context.Context, payload mappings.SignUp) (string, time.Time, error)
	FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error)
	FetchUserPass(ctx context.Context, email string) (string, error)
	FetchUserPassById(ctx context.Context, userId string) (string, error)
	ChangePass(ctx context.Context, userId string, newPass string) error
	GetEmail(ctx context.Context, userId string) (string, error)
	ChangeEmail(ctx context.Context, userId string, newEmail string) error
	GetRespTime(ctx context.Context) *time.Duration
}

type Db struct {
	Auth interface {
		CreateUser(ctx context.Context, payload mappings.SignUp) (string, time.Time, error)
	}
	Password interface {
		FetchUserPass(ctx context.Context, email string) (string, error)
		FetchUserPassById(ctx context.Context, userId string) (string, error)
		ChangePass(ctx context.Context, userId string, newPass string) error
	}
	Email interface {
		GetEmail(ctx context.Context, userId string) (string, error)
		ChangeEmail(ctx context.Context, userId string, newEmail string) error
	}
	UserDeatails interface {
		FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error)
	}
	Metrics interface{
		GetRespTime(ctx context.Context) *time.Duration
	}
}

func NewDb2(pg *pgxpool.Pool) Db {
	return Db{
		Auth: &authDb{pg},
		Password: NewPassDb(pg),
		Email: NewEmailDb(pg),
		UserDeatails: NewUserDb(pg),
		Metrics: NewMetricsDb(pg),
	}
}

