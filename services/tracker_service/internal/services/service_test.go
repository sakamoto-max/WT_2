package services

import (
	"context"
	"testing"
	"tracker_service/internal/client/clientmock"
	"tracker_service/internal/domain"
	"tracker_service/internal/repository"
	"tracker_service/internal/repository/cache"
	"tracker_service/internal/repository/cache/cachemock"
	"tracker_service/internal/repository/mock"

	"github.com/go-openapi/testify/assert"
	trackerpb "github.com/sakamoto-max/wt_2_proto/shared/tracker"
)

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
					Metrics: &mock.MetricsMock{
						PgDown: test.pgErr,
					},
				},
				cache: &cache.Cache{
					Metrics: &cachemock.MetricsMock{
						Down: test.redisErr,
					},
				},
			}

			resp, err := s.GetHealth(context.Background(), &trackerpb.GetHealthReq{})
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

func Test_StartEmptyWorkout(t *testing.T) {

	tests := []struct {
		name           string
		pgDown         bool
		cacheDown      bool
		clientDown     bool
		planExists     bool
		workoutOngoing bool
		wantErr        bool
	}{
		{name: "success", planExists: true, wantErr: false, workoutOngoing: false},
		{name: "pg is down", pgDown: true, wantErr: true, workoutOngoing: false},
		{name: "cache is down", cacheDown: true, wantErr: true, workoutOngoing: false},
		{name: "client is down", clientDown: true, workoutOngoing: false, wantErr: true, planExists: true},
		{name: "workout is alredy ongoing", workoutOngoing: true, wantErr: true},
		{name: "empty plan doesn't exist", planExists: false, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					Start: &mock.StartMock{
						PgDown:         test.pgDown,
						WorkoutOngoing: test.workoutOngoing,
					},
				},
				cache: &cache.Cache{
					TrackerId: &cachemock.TrackerIdMock{
						Down:           test.cacheDown,
						WorkoutOngoing: test.workoutOngoing,
					},
				},
				planClient: &clientmock.PlanClientMock{
					Down:      test.clientDown,
					PlanExits: test.planExists,
				},
			}

			resp, err := s.StartEmptyWorkout(context.TODO(), &trackerpb.StartEmptyWorkoutReq{UserId: "123"})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.Message)
		})
	}
}

func Test_StartWorkoutWithPlan(t *testing.T) {
	// success
	// pg is down
	// cache is down
	// client is down
	// plan doesn't exist
	// workout ongoing alr

	tests := []struct {
		name           string
		pgDown         bool
		cacheDown      bool
		clientDown     bool
		workoutOngoing bool
		planExists     bool
		wantErr        bool
	}{
		{name: "success", wantErr: false, planExists: true, workoutOngoing: false},
		{name: "pg is down", wantErr: true, pgDown: true, workoutOngoing: false},
		{name: "cache is down", wantErr: true, cacheDown: true, workoutOngoing: false},
		{name: "client is down", wantErr: true, clientDown: true, workoutOngoing: false},
		{name: "plan doesn't exist", wantErr: true, workoutOngoing: false, planExists: false},
		{name: "workout ongoing already", wantErr: true, workoutOngoing: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					Start: &mock.StartMock{
						PgDown:         test.pgDown,
						WorkoutOngoing: test.workoutOngoing,
					},
				},

				cache: &cache.Cache{
					TrackerId: &cachemock.TrackerIdMock{
						Down:           test.cacheDown,
						WorkoutOngoing: test.workoutOngoing,
					},

					CurrentPlan: &cachemock.CurrentPlanMock{
						Down:           test.cacheDown,
						WorkoutOngoing: test.workoutOngoing,
					},

					Plan: &cachemock.PlanMock{
						Down:           test.cacheDown,
						WorkoutOngoing: test.workoutOngoing,
					},
				},
				planClient: &clientmock.PlanClientMock{
					Down:      test.clientDown,
					PlanExits: test.planExists,
				},
			}

			resp, err := s.StartWorkoutWithPlan(context.Background(), &trackerpb.StartWorkoutWithPlanReq{UserId: "123", PlanName: "planName"})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.ExercisesInPlan)
			assert.NotZero(t, resp.Message)
			assert.NotZero(t, resp.PlanName)

		})
	}
}

func Test_CancelWorkout(t *testing.T) {
	// depends on cache & pg

	// success
	// pg down
	// cache down
	// no ongoing workout

	tests := []struct {
		name           string
		pgDown         bool
		cacheDown      bool
		workoutOngoing bool
		cacheHit       bool
		wantErr        bool
		withPlan       bool
	}{
		{name: "success", workoutOngoing: true, cacheHit: true, wantErr: false},
		{name: "success - with plan", workoutOngoing: true, cacheHit: true, wantErr: false, withPlan: true},
		{name: "pg is down", workoutOngoing: true, pgDown: true, wantErr: true},
		{name: "cache is down", workoutOngoing: true, cacheDown: true, wantErr: true},
		{name: "no ongoing workout", workoutOngoing: false, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					Cancel: &mock.CancelMock{
						PgDown:         test.pgDown,
						WorkoutOngoing: test.workoutOngoing,
					},
				},
				cache: &cache.Cache{
					TrackerId: &cachemock.TrackerIdMock{
						Down:           test.cacheDown,
						WorkoutOngoing: test.workoutOngoing,
					},
					CurrentPlan: &cachemock.CurrentPlanMock{
						Down:           test.cacheDown,
						WorkoutOngoing: test.workoutOngoing,
						WithPlan:       test.withPlan,
					},
					UserData: &cachemock.UserDataMock{
						Down: test.cacheDown,
					},
				},
			}

			resp, err := s.CancelWorkout(context.Background(), &trackerpb.CancelWorkoutReq{UserId: "123"})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.Message)
		})
	}

}

func Test_EndWorkout(t *testing.T) {

	dataNoChange := trackerpb.EndWorkoutReq{
		UserId: "123",
		AllExerices: []*trackerpb.TrackerForEachExer{
			{
				ExerciseName: "exer_1",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
			{
				ExerciseName: "exer_2",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
			{
				ExerciseName: "exer_3",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
		},
	}

	dataNewExers := trackerpb.EndWorkoutReq{
		UserId: "123",
		AllExerices: []*trackerpb.TrackerForEachExer{
			{
				ExerciseName: "exer_1",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
			{
				ExerciseName: "exer_2",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
			{
				ExerciseName: "exer_3",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
			{
				ExerciseName: "exer_4",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
			{
				ExerciseName: "exer_4",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
		},
	}
	dataExerNotPerformedWithYes := trackerpb.EndWorkoutReq{
		UserId:       "123",
		UserResponse: true,
		AllExerices: []*trackerpb.TrackerForEachExer{
			{
				ExerciseName: "exer_1",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
			{
				ExerciseName: "exer_2",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
		},
	}
	dataExerNotPerformedWithNo := trackerpb.EndWorkoutReq{
		UserId:       "123",
		UserResponse: false,
		AllExerices: []*trackerpb.TrackerForEachExer{
			{
				ExerciseName: "exer_1",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
			{
				ExerciseName: "exer_2",
				SetsAndReps: []*trackerpb.SetsAndReps{
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
					{Reps: 10, Weight: 12.5},
				},
			},
		},
	}

	tests := []struct {
		name           string
		pgDown         bool
		cacheDown      bool
		withPlan       bool
		workoutOngoing bool
		conflictLevel  int
		clientDown     bool
		exerExists     bool
		wantErr        bool
		wantResp       bool
		data           *trackerpb.EndWorkoutReq
	}{
		{name: "success - empty workout", wantErr: false, withPlan: false, workoutOngoing: true, exerExists: true, data: &dataNoChange},
		{name: "success - with plan - no changes", wantErr: false, withPlan: true, workoutOngoing: true, exerExists: true, conflictLevel: 0, data: &dataNoChange},
		{name: "some exercises not performed", workoutOngoing: true, withPlan: true, data: &dataExerNotPerformedWithYes, exerExists: true, conflictLevel: 0, wantErr: false, wantResp: true},
		{name: "success - some exercises not performed - user responded with yes", withPlan: true, workoutOngoing: true, data: &dataExerNotPerformedWithYes, exerExists: true, conflictLevel: 1, wantErr: false},
		{name: "success - some exercises not performed - user responded with no", withPlan: true, workoutOngoing: true, data: &dataExerNotPerformedWithNo, exerExists: true, conflictLevel: 1, wantErr: false},
		{name: "new exercises added", withPlan: true, workoutOngoing: true, data: &dataNewExers, wantErr: false, exerExists: true, conflictLevel: 0, wantResp: true},
		{name: "success - new exercises added - user responded with yes", withPlan: true, workoutOngoing: true, data: &dataNewExers, wantErr: false, exerExists: true, conflictLevel: 2},
		{name: "success - new exercises added - user respondes with no", withPlan: true, workoutOngoing: true, data: &dataNewExers, wantErr: false, exerExists: true, conflictLevel: 2},
		{name: "no workout ongoing", workoutOngoing: false, wantErr: true, data: &dataNoChange} ,
		{name: "cache down", wantErr: true, cacheDown: true, data: &dataNewExers},
		{name: "pg down", wantErr: true, pgDown: true, workoutOngoing: true, data: &dataNoChange},
		{name: "client down", wantErr: true, workoutOngoing: true, clientDown: true, data: &dataNoChange},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				pg: &repository.Db{
					End: &mock.EndMock{
						PgDown:         test.pgDown,
						WorkoutOngoing: test.workoutOngoing,
					},
				},
				cache: &cache.Cache{
					TrackerId: &cachemock.TrackerIdMock{
						Down:           test.cacheDown,
						WorkoutOngoing: test.workoutOngoing,
					},
					CurrentPlan: &cachemock.CurrentPlanMock{
						Down:           test.cacheDown,
						WorkoutOngoing: test.workoutOngoing,
						WithPlan:       test.withPlan,
					},
					Plan: &cachemock.PlanMock{
						Down:           test.cacheDown,
						WorkoutOngoing: test.workoutOngoing,
					},
					Conflict: &cachemock.ConflictMock{
						Down:          test.cacheDown,
						ConflictLevel: test.conflictLevel,
					},
					TrackerData: &cachemock.TrackerDataMock{
						Down: test.cacheDown,
						Data: domain.ConvertToLocal(test.data),
					},
					NewExercises: &cachemock.NewExercisesMock{
						Down: test.cacheDown,
					},
					UserData: &cachemock.UserDataMock{
						Down: test.cacheDown,
					},
				},
				exerClient: &clientmock.ExerClientMock{
					Down:           test.clientDown,
					ExerciseExists: test.exerExists,
				},
			}

			resp, err := s.EndWorkout(context.Background(), test.data)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			if test.wantResp {
				assert.NoError(t, err)
				assert.True(t, resp.ConflictOccured)
				assert.NotEmpty(t, resp.ExerciseNames)
				assert.NotZero(t, resp.Message)
				assert.NotZero(t, resp.Reason)
				assert.NotZero(t, resp.RequestStatus)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.Message)
		})
	}

}
