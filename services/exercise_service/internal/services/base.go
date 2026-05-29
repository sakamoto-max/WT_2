package services

import (
	"exercise_service/internal/repository"
	"exercise_service/internal/repository/cache"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
)

func NewService(pg *repository.Db, cache *cache.Cache) *Service {
	return &Service{
		pg:    pg,
		cache: cache,
	}
}
type Service struct {
	exerpb.UnimplementedExerciseServiceServer
	pg    *repository.Db
	cache *cache.Cache
}


