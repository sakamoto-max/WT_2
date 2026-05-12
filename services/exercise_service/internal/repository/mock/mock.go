package mock

import (
	"context"
	"exercise_service/internal/domain"
	"exercise_service/internal/repository"
	"time"
)

type mockDb struct{}

func NewMockDb() repository.RepoIface {
	return &mockDb{}
}

var (
	timeDuration time.Duration = time.Second
	exerciseName string = "mock_exer"
	exerId string = "b75f876d-e916-40af-8f97-1c4f783fb40c"
	restTime int = 180
    bodyPart  string  = "mock_body_part"
    equipment string  = "mock_equipment" 
	createdAt time.Time = time.Now()
    updatedAt time.Time = time.Now()
)

func (m *mockDb) GetPostgresRespTime(ctx context.Context) *time.Duration {
	return &timeDuration
}
func (m *mockDb) GetRedisRespTime(ctx context.Context) *time.Duration    {
	return &timeDuration
}
func (m *mockDb) GetExerciseByNameR(ctx context.Context, userId string, exerciseName string) (*domain.Exercise, error) {

	return &domain.Exercise{
		Id: exerId,
		Name: exerciseName,
		RestTime: restTime,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
		BodyPart: bodyPart,
		Equipment: equipment,
	}, nil
}
func (m *mockDb) GetExerciseByName(ctx context.Context, userId string, exerciseName string) (*domain.Exercise, error) {
	return &domain.Exercise{
		Id: exerId,
		Name: exerciseName,
		RestTime: restTime,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
		BodyPart: bodyPart,
		Equipment: equipment,
	}, nil
}
func (m *mockDb) SetExerciseByNameR(ctx context.Context, userId string, exerData *domain.Exercise) error {
	return nil
}
func (m *mockDb) CreateExercise(ctx context.Context, userId string, exerciseName string, bodyPartName string, equipmentName string) (string, error) {
	return exerId, nil
}
func (m *mockDb) GetAllExercises(ctx context.Context, userId string) (*[]domain.Exercise, error) {
	exer := domain.Exercise{
		Id: exerId,
		Name: exerciseName,
		RestTime: restTime,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
		BodyPart: bodyPart,
		Equipment: equipment,
	}
	exers := []domain.Exercise{exer, exer, exer, exer}

	return &exers, nil
}
func (m *mockDb) GetExerciseNameByID(ctx context.Context, exerciseId string) (string, error)     {
	return exerciseName, nil
}
func (m *mockDb) DeleteExecise(ctx context.Context, userId string, exerciseName string) error    {
	return nil
}
func (m *mockDb) ExerciseExistsReturnId(ctx context.Context, userId string, exerciseName string) (string, error) {
	return exerId, nil
}
