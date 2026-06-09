package services

import (
	"context"
	"exercise_service/internal/mappings"
	"fmt"

	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) GetOneExercise(ctx context.Context, in *exerpb.SendExerciseName) (*exerpb.OneExerciseResp, error) {
	Exercise, err := s.cache.CRUD.GetExerciseByNameR(ctx, mappings.ToGetExerciesByName(in))
	if err != nil || Exercise == nil {

		Exercise, err = s.pg.ExerciseGet.GetExerciseByName(ctx, mappings.ToGetExerciesByName(in))
		if err != nil {
			return nil, err
		}

		s.cache.CRUD.SetExerciseByNameR(ctx, in.UserId, Exercise)
	}

	fmt.Println()

	return &exerpb.OneExerciseResp{
		Id:        Exercise.Id,
		Name:      Exercise.Name,
		BodyPart:  Exercise.BodyPart,
		Equipment: Exercise.Equipment,
		CreatedAt: timestamppb.New(Exercise.CreatedAt),
		UpdatedAt: timestamppb.New(Exercise.UpdatedAt),
	}, nil
}
func (s *Service) GetAllExercises(ctx context.Context, in *exerpb.GetAllExercisesREq) (*exerpb.GetAllExercisesResp, error) {



	allExercises, err := s.pg.ExerciseGet.GetAllExercises(ctx, mappings.ToGetAllExercises(in))
	if err != nil {
		return nil, err
	}

	resp := exerpb.GetAllExercisesResp{}

	for _, exer := range *allExercises {
		eachExer := exerpb.OneExerciseResp{
			Id:        exer.Id,
			Name:      exer.Name,
			BodyPart:  exer.BodyPart,
			Equipment: exer.Equipment,
			CreatedAt: timestamppb.New(exer.CreatedAt),
			UpdatedAt: timestamppb.New(exer.UpdatedAt),
		}

		resp.AllExericses = append(resp.AllExericses, &eachExer)
	}

	resp.NumberOfExercises = int64(len(*allExercises))

	return &resp, nil
}
func (s *Service) CreateExercise(ctx context.Context, in *exerpb.CreateExerciseReq) (*exerpb.CreateExerciseResp, error) {
	UUId, err := s.pg.ExerciseCD.CreateExercise(ctx, mappings.ToCreateExercise(in))
	if err != nil {
		return nil, err
	}

	return &exerpb.CreateExerciseResp{
		Id: UUId,
	}, nil
}
func (s *Service) DeleteExercise(ctx context.Context, in *exerpb.SendExerciseName) (*exerpb.DeleteExerciseResp, error) {

	err := s.pg.ExerciseCD.DeleteExecise(ctx, mappings.ToDeleteExercise(in))
	if err != nil {
		return nil, err
	}

	err = s.cache.CRUD.DeleteExerciseByNameR(ctx, mappings.ToDeleteExercise(in))
	if err != nil {
		return nil, err
	}

	return &exerpb.DeleteExerciseResp{
		Message: fmt.Sprintf("exercise %s is deleted", in.ExerciseName),
	}, nil
}
func (s *Service) ExerciseExistsReturnId(ctx context.Context, in *exerpb.SendExerciseName) (*exerpb.ExerciseExistsReturnIdResp, error) {
	exerciseId, err := s.pg.ExerciseGet.ExerciseExistsReturnId(ctx, mappings.ToExerciseExistsReturnId(in))
	if err != nil {
		return nil, err
	}

	return &exerpb.ExerciseExistsReturnIdResp{
		ExerciseId: exerciseId,
	}, nil
}
func (s *Service) GetExerciseName(ctx context.Context, in *exerpb.SendExerciseID) (*exerpb.GetExerciseNameResp, error) {
	exerciseName, err := s.pg.ExerciseGet.GetExerciseNameByID(ctx, in.ExerciseId)
	if err != nil {
		return nil, err
	}

	return &exerpb.GetExerciseNameResp{
		ExerciseName: exerciseName,
	}, nil
}
func (s *Service) PING(ctx context.Context, in *exerpb.PingExerReq) (*exerpb.PingExerResp, error) {
	return &exerpb.PingExerResp{}, nil
}
func (s *Service) GetHealth(ctx context.Context, in *exerpb.GetHealthReq) (*exerpb.GetHealthResp, error) {
	pgRespTime := s.pg.Metrics.GetRespTime(ctx)
	redisRespTime := s.cache.Metrics.GetRespTime(ctx)

	if pgRespTime == nil && redisRespTime == nil {
		return &exerpb.GetHealthResp{
			PostgresRespTime: nil,
			RedisRespTime:    nil,
		}, nil
	}
	if redisRespTime == nil {
		return &exerpb.GetHealthResp{
			PostgresRespTime: durationpb.New(*pgRespTime),
			RedisRespTime:    nil,
		}, nil
	}
	if pgRespTime == nil {
		return &exerpb.GetHealthResp{
			PostgresRespTime: nil,
			RedisRespTime:    durationpb.New(*redisRespTime),
		}, nil
	}

	return &exerpb.GetHealthResp{
		PostgresRespTime: durationpb.New(*pgRespTime),
		RedisRespTime:    durationpb.New(*redisRespTime),
	}, nil
}
