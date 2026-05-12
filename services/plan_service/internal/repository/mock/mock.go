package mock

import (
	"context"
	"plan_service/internal/models"
	"plan_service/internal/repository"
	"time"
)

var (
	emptyPlanId string = "b343fa55-9e2a-4bfe-9be0-5fb5f3dbdc4d"
	planId string = "3af5d6d0-8be6-4a39-84f7-e20323a3e4a1"
	timeDuration time.Duration = time.Second
	exerciseId string = "9d497c95-8913-4c70-8ea0-6891bbd17af5"
)


type mockDb struct{}

func NewMockDb() repository.RepoIFace {
	return &mockDb{}
}

func (m *mockDb) CreatePlan(ctx context.Context, userId string, planName string, exerciseIds []string) error {
	return nil
}
func (m *mockDb) GetPlans(ctx context.Context, userId string) (*[]models.Plan3, error) {

	resp := []models.Plan3{
		{"mock_plan_1", planId},
		{"mock_plan_2", planId},
		{"mock_plan_3", planId},
		{"mock_plan_4", planId},
		{"mock_plan_5", planId},
		{"mock_plan_6", planId},
	}

	return &resp, nil
}


func (m *mockDb) GetAllExercisesByPlanID(ctx context.Context, planId string) (*[]string, error)     {
	return &[]string{exerciseId, exerciseId, exerciseId, exerciseId}, nil
}
func (m *mockDb) ReturnsPlanId(ctx context.Context, userId string, planName string) (string, error) {
	return planId, nil
}
func (m *mockDb) AddExercisesToPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {
	return nil
}
func (m *mockDb) DeleteExerciseFromPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {
	return nil
}
func (m *mockDb) DeletePlan(ctx context.Context, userId string, planId string) error          {
	return nil
}

func (m *mockDb) CreateEmptyPlan(ctx context.Context, userId string) error                    {
	return nil
}
func (m *mockDb) GetPostgresRespTime(ctx context.Context) *time.Duration                      {
	return &timeDuration

}
func (m *mockDb) GetRedisRespTime(ctx context.Context) *time.Duration                         {
	return &timeDuration
}
func (m *mockDb) SetUserEmptyPlanIdR(ctx context.Context, userId string, planId string) error {
	return nil
}
func (m *mockDb) GetUserEmptyPlanIdR(ctx context.Context, userId string) (string, error)      {
	return emptyPlanId, nil
}
func (m *mockDb) DelUserEmptyPlanIdR(ctx context.Context, userId string) error                {
	return nil
}
func (m *mockDb) Close() error                                                                {
	return nil
}
func (m *mockDb) GetPlan(ctx context.Context, userId string, planName string) (string, *[]string, error) {
	return planId, &[]string{exerciseId, exerciseId, exerciseId}, nil

}
func (m *mockDb) SetUserPlanId(ctx context.Context, userId string, planName string, planId string) error {
	return nil
}
func (m *mockDb) GetUserPlanId(ctx context.Context, userId string, planName string) (string, error) {
	return planId, nil
}

func (m *mockDb) DelUserPlanId(ctx context.Context, userId string, planName string) error           {
	return nil
}
func (m *mockDb) SetUserPlan(ctx context.Context, userId string, planName string, planId string, exerciseIds *[]string) error {
	return nil
}
func (m *mockDb) GetUserPlan(ctx context.Context, userId string, planName string) (string, *[]string, error) {
	return planId, &[]string{exerciseId, exerciseId, exerciseId}, nil

}
func (m *mockDb) DelUserPlan(ctx context.Context, userId string, planName string) error {
	return nil
}
