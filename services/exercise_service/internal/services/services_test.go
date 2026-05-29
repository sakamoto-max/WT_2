package services

import (
	"context"
	"exercise_service/internal/repository"
	"exercise_service/internal/repository/cache"
	"exercise_service/internal/repository/cache/cachemock"
	"exercise_service/internal/repository/mock"
	"testing"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"github.com/stretchr/testify/assert"
)

func Test_GetExerciseByName(t *testing.T) {

	tests := []struct {
		name      string
		DbErr     bool
		CacheDown bool
		ExerExits bool
		Hit       bool
		Miss      bool
		WantErr   bool
	}{
		{name: "success - cache miss", Miss: true, ExerExits: true},
		{name: "success - cache hit", WantErr: false, Hit: true},
		{name: "success - cache down", WantErr: false, CacheDown: true, ExerExits: true},
		{name: "exercise doesn't exist", WantErr: true, ExerExits: false, Miss: true, DbErr: false, CacheDown: false},
		{name: "pg has error", WantErr: true, Miss: true, DbErr: true},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					ExerciseGet: &mock.ExerciseGetMock{
						HasErr:        test.DbErr,
						ExerciseExits: test.ExerExits,
					},
				},
				cache: &cache.Cache{
					CRUD: &cachemock.CrudMock{
						Down: test.CacheDown,
						Miss: test.Miss,
						Hit:  test.Hit,
					},
				},
			}

			resp, err := s.GetOneExercise(context.Background(), &exerpb.SendExerciseName{UserId: "123", ExerciseName: "exer"})
			if test.WantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.BodyPart)
			assert.NotZero(t, resp.CreatedAt)
			assert.NotZero(t, resp.Equipment)
			assert.NotZero(t, resp.Id)
			assert.NotZero(t, resp.Name)
			assert.NotZero(t, resp.UpdatedAt)
		})
	}
}
func Test_GetAllExercisesSer(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
		PgDown  bool
	}{
		{name: "success", wantErr: false, PgDown: false},
		{name: "pg down", wantErr: true, PgDown: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					ExerciseGet: &mock.ExerciseGetMock{
						HasErr: test.PgDown,
					},
				},
			}

			allExers, err := s.GetAllExercises(context.Background(), &exerpb.GetAllExercisesREq{})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, allExers.AllExericses)
		})
	}
}
func Test_DeleteExerciseSet(t *testing.T) {

	tests := []struct {
		name      string
		wantErr   bool
		PgDown    bool
		ExerExits bool
		cacheDown bool
	}{
		{name: "success", wantErr: false, PgDown: false, ExerExits: true},
		{name: "exercise doesn't exist", wantErr: true, PgDown: false, ExerExits: false},
		{name: "pg down", wantErr: true, PgDown: true},
		{name: "cache down", wantErr: true, PgDown: false, cacheDown: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				pg: &repository.Db{
					ExerciseCD: &mock.ExerciseCdMock{
						HasErr:         test.PgDown,
						ExerciseExists: test.ExerExits,
					},
				},
				cache: &cache.Cache{
					CRUD: &cachemock.CrudMock{
						Down: test.cacheDown,
					},
				},
			}

			_, err := s.DeleteExercise(context.Background(), &exerpb.SendExerciseName{UserId: "123", ExerciseName: "exer"})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			// assert.(t, resp.Message)
		})
	}

}
func Test_CreateExerciseSer(t *testing.T) {

	tests := []struct {
		name       string
		pgDown     bool
		wantErr    bool
		exerExists bool
	}{
		{name: "success", pgDown: false, wantErr: false, exerExists: false},
		{name: "pg down", pgDown: true, wantErr: true, exerExists: false},
		{name: "exercise already exits", pgDown: false, wantErr: true, exerExists: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				pg: &repository.Db{
					ExerciseCD: &mock.ExerciseCdMock{
						HasErr:         test.pgDown,
						ExerciseExists: test.exerExists,
					},
				},
			}

			resp, err := s.CreateExercise(context.Background(), &exerpb.CreateExerciseReq{})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.Id)
		})
	}

}
func Test_ExerciseExistsReturnId(t *testing.T) {

	tests := []struct {
		name           string
		exerciseExists bool
		wantErr        bool
		pgDown         bool
	}{
		{name: "exercise exists", exerciseExists: true, wantErr: false},
		{name: "exercise doesn't exist", exerciseExists: false, wantErr: true},
		{name: "pg is down", exerciseExists: true, wantErr: true, pgDown: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				pg: &repository.Db{
					ExerciseGet: &mock.ExerciseGetMock{
						HasErr:        test.pgDown,
						ExerciseExits: test.exerciseExists,
					},
				},
			}

			resp, err := s.ExerciseExistsReturnId(context.Background(), &exerpb.SendExerciseName{})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			
			assert.NoError(t, err)
			assert.NotZero(t, resp.ExerciseId)
		})

	}
}
func Test_GetExerciseNameByID(t *testing.T) {

	tests := []struct {
		name           string
		exerciseExists bool
		wantErr        bool
		pgDown         bool
	}{
		{name: "exercise exists", exerciseExists: true, wantErr: false},
		{name: "exercise doesn't exist", exerciseExists: false, wantErr: true},
		{name: "pg is down", exerciseExists: true, wantErr: true, pgDown: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				pg: &repository.Db{
					ExerciseGet: &mock.ExerciseGetMock{
						HasErr:        test.pgDown,
						ExerciseExits: test.exerciseExists,
					},
				},
			}

			resp, err := s.GetExerciseName(context.Background(), &exerpb.SendExerciseID{})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.ExerciseName)
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

		s := Service{
			pg: &repository.Db{
				Metrics: &mock.MetricsDbMock{
					HasErr: test.pgErr,
				},
			},

			cache: &cache.Cache{
				Metrics: &cachemock.MetricsMock{
					Down: test.redisErr,
				},
			},
		}

		resp, err := s.GetHealth(context.Background(), &exerpb.GetHealthReq{})
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

	}

}
