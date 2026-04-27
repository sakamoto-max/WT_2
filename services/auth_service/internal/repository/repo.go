package repository

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type repo struct {
	pDB *pgxpool.Pool
	rDB *redis.Client
}

type RepoIface interface {
	GetEmail(ctx context.Context, userId string) (string, error)
	ChangeEmail(ctx context.Context, userId string, newEmail string) error
	GetRefreshToken(ctx context.Context, uuid string) (string, error)
	SetRefreshTokenAndUUID(ctx context.Context, uuid string, Refreshtoken string, userId string) error
	SetUUID(ctx context.Context, uuid string, userId string) error
	GetUUID(ctx context.Context, userId string) (string, error)
	RefreshExists(ctx context.Context, userId string) (bool, error)
	CreateUser(ctx context.Context, name string, email string, hashedPass string, role string) (string, time.Time, error)
	FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error)
	UserLogout(ctx context.Context, userId string, uuid string) error
	FetchUserPass(ctx context.Context, email string) (string, error)
	ChangePass(ctx context.Context, userId string, newPass string) error
	FetchUserPassById(ctx context.Context, userId string) (string, error)
	GetPostgresRespTime(ctx context.Context) *time.Duration
	GetRedisRespTime(ctx context.Context) *time.Duration
}

func NewRepo(pool *pgxpool.Pool, redisClient *redis.Client) RepoIface {
	return &repo{pDB: pool, rDB: redisClient}
}

func (r *repo) GetPostgresRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.pDB.Ping(ctx)
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}

func (r *repo) GetRedisRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.rDB.Ping(ctx).Err()
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
