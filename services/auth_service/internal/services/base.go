package services

import (
	"auth_service/internal/repository"
	"auth_service/internal/repository/cache"
	"errors"
	pb "github.com/sakamoto-max/wt_2_proto/shared/auth"
)

var (
	ErrEmailDoesntMatch  = errors.New("email user sent is wrong")
	ErrPasswordIncorrect = errors.New("incorrect password")
	ErrSameOldPass       = errors.New("new password cannot be same as the old password")
	ErrSameOldEmail      = errors.New("new email cannot be same as the old email")
)

func NewService(pg repository.Db, client *cache.Cache) *Service {
	return &Service{
		db:    pg,
		cache: client,
	}
}

type Service struct {
	pb.UnimplementedAuthServiceServer
	db    repository.Db
	cache *cache.Cache
}
