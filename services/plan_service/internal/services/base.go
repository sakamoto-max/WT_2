package services

import (
	"context"
	"plan_service/internal/client"
	"plan_service/internal/repository"
	"plan_service/internal/repository/cache"
	// "strings"
	"google.golang.org/protobuf/types/known/durationpb"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	// exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
)

type Service struct {
	pg      *repository.Db
	cache   *cache.Cache
	gClient client.ExerClientIface
	planpb.UnimplementedPlanServiceServer
}

func NewService(Db *repository.Db, cache *cache.Cache, grpcCli client.ExerClientIface) *Service {
	return &Service{pg: Db, cache: cache, gClient: grpcCli}
}


func (s *Service) GetHealth(ctx context.Context, in *planpb.GetHealthReq) (*planpb.GetHealthResp, error) {
	pgRespTime := s.pg.MetricsRepo.GetRespTime(ctx)
	redisRespTime := s.cache.Metrics.GetRespTime(ctx)

	if pgRespTime == nil && redisRespTime == nil {
		return &planpb.GetHealthResp{
			PostgresRespTime: nil,
			RedisRespTime:    nil,
		}, nil
	}
	if redisRespTime == nil {
		return &planpb.GetHealthResp{
			PostgresRespTime: durationpb.New(*pgRespTime),
			RedisRespTime:    nil,
		}, nil
	}
	if pgRespTime == nil {
		return &planpb.GetHealthResp{
			PostgresRespTime: nil,
			RedisRespTime:    durationpb.New(*redisRespTime),
		}, nil
	}

	return &planpb.GetHealthResp{
		PostgresRespTime: durationpb.New(*pgRespTime),
		RedisRespTime:    durationpb.New(*redisRespTime),
	}, nil
}
