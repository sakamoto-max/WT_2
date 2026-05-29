package services

import (
	"tracker_service/internal/client"
	"tracker_service/internal/repository"
	"tracker_service/internal/repository/cache"
	trackerpb "github.com/sakamoto-max/wt_2_proto/shared/tracker"
)

type Service struct {
	trackerpb.UnimplementedTrackerServiceServer
	pg         *repository.Db
	cache      *cache.Cache
	planClient client.PlanClientIFace
	exerClient client.ExerClientIface
}



func NewService(Db *repository.Db, cache *cache.Cache, planClient client.PlanClientIFace, exerClient client.ExerClientIface) *Service {
	return &Service{pg: Db, cache: cache, planClient: planClient, exerClient: exerClient}
}