package services

import (
	"context"
	"plan_service/internal/client/clientmock"

	// "plan_service/internal/domain/plan"

	// "plan_service/internal/domain/plan"

	// "plan_service/internal/domain/plan"
	"plan_service/internal/repository"
	"plan_service/internal/repository/cache"
	"plan_service/internal/repository/cache/cachemock"
	"plan_service/internal/repository/mock"
	"testing"

	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	"github.com/stretchr/testify/assert"
)

func Test_CreatePlan(t *testing.T) {

	tests := []struct {
		name           string
		pgDown         bool
		grpcServerDown bool
		planExists     bool
		exerciseExists bool
		wantErr        bool
	}{
		{name: "success", wantErr: false, exerciseExists: true},
		{name: "plan alr exits", planExists: true, wantErr: true},
		{name: "pg is down", wantErr: true, exerciseExists: true, pgDown: true},
		{name: "grpc server down", wantErr: true, exerciseExists: true, grpcServerDown: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					PlanCommandRepo: &mock.PlanCommandMock{
						PgDown:     test.pgDown,
						PlanExists: test.planExists,
					},
				},
				gClient: &clientmock.ClientMock{
					ServerIsDown:   test.grpcServerDown,
					ExerciseExists: test.exerciseExists,
				},
			}

			resp, err := s.CreatePlan(context.Background(), &planpb.CreatePlanReq{PlanName: "plan", ExerciseNames: []string{"abc", "def"}})
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.ExerciseNames)
			assert.NotZero(t, resp.Message)
			assert.NotZero(t, resp.PlanName)
		})

	}

}

func Test_GetPlans(t *testing.T) {
	// success
	// pg is down
	// client is down
	// no plan exists

	tests := []struct {
		name           string
		pgDown         bool
		clientDown     bool
		wantErr        bool
		planExists     bool
		exerciseExists bool
	}{
		{name: "success", wantErr: false, planExists: true, exerciseExists: true},
		{name: "pg is down", wantErr: true, pgDown: true, planExists: true, exerciseExists: true},
		{name: "client is down", wantErr: true, clientDown: true, planExists: true, exerciseExists: true},
		{name: "no plans exists", wantErr: true, planExists: false, exerciseExists: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					PlanQueryRepo: &mock.PlanQueryMock{
						PgDown:    test.pgDown,
						PlanExits: test.planExists,
					},
				},
				gClient: &clientmock.ClientMock{
					ServerIsDown:   test.clientDown,
					ExerciseExists: test.exerciseExists,
				},
			}

			resp, err := s.GetAllPlans(context.Background(), &planpb.GetAllPlansReq{UserId: "123"})
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.AllPlans)
			assert.NotZero(t, resp.NumberOfPlans)
		})
	}
}

func Test_GetPlan(t *testing.T) {
	// depends on cache, pg, client

	// success
	// cache miss
	// plan doesn't exist -> cache miss -> plan doesn't exist in pg
	// pg is down ->
	// -> cache hit
	// -> cache miss
	// cache is down ->
	// plan exists
	// plan doesn't exist
	// grpc server is down

	tests := []struct {
		name           string
		planExists     bool
		pgDown         bool
		CacheDown      bool
		CacheHit       bool
		grpcDown       bool
		exerciseExists bool
		wantErr        bool
	}{
		{name: "success", CacheHit: true, wantErr: false, exerciseExists: true},
		{name: "success - cache miss", planExists: true, wantErr: false, exerciseExists: true},
		{name: "plan doesn't exist", planExists: false, wantErr: true},
		{name: "pg down - plan exists in cache", planExists: true, CacheHit: true, wantErr: false, exerciseExists: true},
		{name: "pg down - plan doesn't exist in cache", planExists: false, wantErr: true},
		{name: "cache down - plan exists", wantErr: false, planExists: true, exerciseExists: true, CacheDown: true},
		{name: "cache down - plan doesn't exist", wantErr: true, planExists: false, CacheDown: true},
		{name: "grpc server down", wantErr: true, planExists: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				pg: &repository.Db{
					PlanQueryRepo: &mock.PlanQueryMock{
						PgDown:    test.pgDown,
						PlanExits: test.planExists,
					},
				},
				cache: &cache.Cache{
					UserPlan: &cachemock.UserPlan{
						Down: test.CacheDown,
						Hit:  test.CacheHit,
					},
				},
				gClient: &clientmock.ClientMock{
					ServerIsDown:   test.grpcDown,
					ExerciseExists: test.exerciseExists,
				},
			}

			resp, err := s.GetPlanByName(context.Background(), &planpb.GetPlanByNameReq{UserId: "123", PlanName: "plan"})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotZero(t, resp.PlanId)
			assert.NotZero(t, resp.PlanName)
			assert.NotEmpty(t, resp.ExerciseNames)
		})
	}
}

func Test_AddExercises(t *testing.T) {
	// success
	// success - cache miss
	// pg is down
	// cache - hit
	// cache - miss
	// cache is down
	// pg has plan
	// plan doesnt exist
	// client is down
	// exercise doesn't exist

	tests := []struct {
		name           string
		pgDown         bool
		cacheDown      bool
		clientDown     bool
		planExists     bool
		exerciseExists bool
		cacheHit       bool
		wantErr        bool
	}{
		{name: "success", cacheHit: true, exerciseExists: true, wantErr: false, planExists: true},
		{name: "success - cache miss", cacheHit: false, wantErr: false, planExists: true, exerciseExists: true},
		{name: "pg down - cache hit", wantErr: true, cacheHit: true, exerciseExists: true, pgDown: true},
		{name: "pg down - cache miss", wantErr: true, pgDown: true},
		{name: "cache is down - pg has plan", wantErr: false, planExists: true, exerciseExists: true, cacheDown: true},
		{name: "cache is down - plan doesn't exist", wantErr: true, planExists: false, cacheDown: true},
		{name: "client is down", wantErr: true, planExists: true, clientDown: true},
		{name: "exercise doesn't exists", wantErr: true, exerciseExists: false, planExists: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				pg: &repository.Db{
					PlanQueryRepo: &mock.PlanQueryMock{
						PgDown:    test.pgDown,
						PlanExits: test.planExists,
					},
					PlanExericseRepo: &mock.PlanExericseMock{
						PgDown: test.pgDown,
					},
				},

				cache: &cache.Cache{
					PlanId: &cachemock.PlanId{
						Down: test.cacheDown,
						Hit:  test.cacheHit,
					},
					UserPlan: &cachemock.UserPlan{
						Down: test.cacheDown,
						Hit: test.cacheHit,
					},
				},
				gClient: &clientmock.ClientMock{
					ServerIsDown:   test.clientDown,
					ExerciseExists: test.exerciseExists,
				},
			}

			updatedPlan, err := s.AddExercisesToPlan(context.Background(), &planpb.PlanReq{UserId: "123", PlanName: "plan", ExerciseNames: []string{"exer_1", "exer_2", "exer_3"}})
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, updatedPlan)
			assert.NotZero(t, updatedPlan.PlanName)
			assert.NotEmpty(t, updatedPlan.ExerciseNames)

		})

	}

}

func Test_DeleteExerciseFromPlan(t *testing.T) {

	// success
	// success - cache miss
	// pg is down - err
	// cache is down
	// plan exists in pg
	// plan doesn't exist
	// client is down

	tests := []struct {
		name       string
		pgDown     bool
		cacheDown  bool
		clientDown bool
		planExists bool
		exerExists bool
		wantErr    bool
		cacheHit   bool
	}{
		{name: "success", cacheHit: true, wantErr: false, planExists: true, exerExists: true},
		{name: "success - cache miss", cacheHit: false, wantErr: false, planExists: true, exerExists: true},
		{name: "pg is down", wantErr: true, pgDown: true},
		{name: "cache is down", wantErr: false, cacheDown: true, planExists: true, exerExists: true},
		{name: "cache is down - plan doesn't exist", wantErr: true, cacheDown: true, planExists: false},
		{name: "client is down", wantErr: true, planExists: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				pg: &repository.Db{
					PlanQueryRepo: &mock.PlanQueryMock{
						PgDown:    test.pgDown,
						PlanExits: test.planExists,
					},
					PlanExericseRepo: &mock.PlanExericseMock{
						PgDown: test.pgDown,
					},
				},

				cache: &cache.Cache{
					PlanId: &cachemock.PlanId{
						Down: test.cacheDown,
						Hit:  test.cacheHit,
					},
					UserPlan: &cachemock.UserPlan{
						Down: test.cacheDown,
						Hit: test.cacheHit,
					},
				},
				gClient: &clientmock.ClientMock{
					ServerIsDown:   test.clientDown,
					ExerciseExists: test.exerExists,
				},
			}

			updatedPlan, err := s.DeleteExercisesFromPlan(context.Background(), &planpb.PlanReq{UserId: "123", PlanName: "plan", ExerciseNames: []string{"exer_1", "exer_2", "exer_3"}})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, updatedPlan)
			assert.NotZero(t, updatedPlan.PlanName)
			assert.NotEmpty(t, updatedPlan.ExerciseNames)
		})

	}
}

func Test_DeletePlan(t *testing.T) {

	// success
	// pg is down
	// cache is down
	// cache miss
	// plan doesn't exists

	tests := []struct {
		name       string
		pgDown     bool
		cacheDown  bool
		planExists bool
		cacheHit   bool
		wantErr    bool
	}{
		{name: "success", planExists: true, cacheHit: true, wantErr: false},
		{name: "success - cache miss", wantErr: false, planExists: true, cacheHit: false},
		{name: "pg is down", wantErr: true, pgDown: true},
		{name: "cache is down - pg has plan", wantErr: false, planExists: true, cacheDown: true},
		{name: "cache is down - pg doesn't have plan", wantErr: true, planExists: false, cacheDown: true},
		{name: "plan doesnt exists", wantErr: true, planExists: false, cacheHit: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				pg: &repository.Db{
					PlanQueryRepo: &mock.PlanQueryMock{
						PgDown:    test.pgDown,
						PlanExits: test.planExists,
					},

					PlanCommandRepo: &mock.PlanCommandMock{
						PgDown:     test.pgDown,
						PlanExists: test.planExists,
					},
				},
				cache: &cache.Cache{
					PlanId: &cachemock.PlanId{
						Down: test.cacheDown,
						Hit:  test.cacheHit,
					},

					UserPlan: &cachemock.UserPlan{
						Down: test.cacheDown,
						Hit:  test.cacheHit,
					},
				},
			}

			resp, err := s.DeletePlan(context.TODO(), &planpb.DeletePlanReq{UserId: "123", PlanName: "plan"})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)

		})
	}

}

func Test_GetEmptyPlanId(t *testing.T) {
	// depends on cache, pg

	// success
	// empty plan doesn't exists
	// cache down
	// pg down
	// cache miss

	tests := []struct {
		name       string
		pgDown     bool
		cacheDown  bool
		PlanExists bool
		wantErr    bool
		cacheHit   bool
	}{
		{name: "success", cacheHit: true, wantErr: false},
		{name: "success - cache miss", cacheHit: false, PlanExists: true, wantErr: false},
		{name: "plan doesn't exists", cacheHit: false, PlanExists: false, wantErr: true},
		{name: "pg down - cache hit", cacheHit: true, wantErr: false, pgDown: true},
		{name: "pg down - cache miss", wantErr: true, pgDown: true, PlanExists: true},
		{name: "cache down - pg has data", wantErr: false, cacheDown: true, PlanExists: true},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					PlanQueryRepo: &mock.PlanQueryMock{
						PgDown:    test.pgDown,
						PlanExits: test.PlanExists,
					},
				},
				cache: &cache.Cache{
					EmptyPlan: &cachemock.EmptyPlan{
						Down: test.cacheDown,
						Hit:  test.cacheHit,
					},
				},
			}

			planId, err := s.GetEmptyPlanId(context.Background(), &planpb.SendUserID{UserId: "123"})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, planId)
			assert.NotZero(t, planId)
		})

	}
}

func Test_GetHealth(t *testing.T) {

	tests := []struct {
		name     string
		pgErr    bool
		redisErr bool
		pgNil    bool
		redisNil bool
	}{
		{name: "both working well", pgErr: false, redisErr: false, pgNil: false, redisNil: false},
		{name: "pg is working well and redis is down", pgErr: false, redisErr: true, pgNil: false, redisNil: true},
		{name: "redis is working well and pg is down", pgErr: true, redisErr: false, pgNil: true, redisNil: false},
		{name: "both are down", pgErr: true, redisErr: true, pgNil: true, redisNil: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					MetricsRepo: &mock.MetricsMock{
						PgDown: test.pgErr,
					},
				},
				cache: &cache.Cache{
					Metrics: &cachemock.MetricsMock{
						Down: test.redisErr,
					},
				},
			}

			resp, err := s.GetHealth(context.Background(), &planpb.GetHealthReq{})

			assert.NoError(t, err)

			if test.pgNil && test.redisNil {
				assert.Nil(t, resp.PostgresRespTime)
				assert.Nil(t, resp.RedisRespTime)
				return
			}

			if test.redisNil {
				assert.Nil(t, resp.RedisRespTime)
				assert.NotNil(t, resp.PostgresRespTime)
				return
			}

			if test.pgNil {
				assert.Nil(t, resp.PostgresRespTime)
				assert.NotNil(t, resp.RedisRespTime)
				return
			}

			assert.NotNil(t, resp.PostgresRespTime)
			assert.NotNil(t, resp.RedisRespTime)
				

		})

	}

}
