package mockdb

import (
	"context"
	"fmt"
	"orchestration_service/internal/types"
	"sync"

	"github.com/sakamoto-max/wt_2_proto/shared/enum"
)

type MockDb struct {
	Down    bool
	DbName  string
	HasData bool
	Update  []string
}

func (m *MockDb) FetchData(ctx context.Context, wg *sync.WaitGroup, dataChan chan<- *[]types.Data) {
	defer wg.Done()
	if m.Down {

		data := types.Data{
			Err:         fmt.Errorf("db is down"),
			ServiceName: m.DbName,
		}

		dataChan <- &[]types.Data{data}
		return
	}

	if m.HasData {

		switch m.DbName {
		case enum.ServiceName_AUTH_SERVICE.String():
			authData := types.Data{
				DbId:          "123",
				TargetService: enum.ServiceName_EMAIL_SERVICE.String(),
				CreatedBy:     enum.ServiceName_AUTH_SERVICE.String(),
				Task:          enum.TaskName_SEND_EMAIL_FOR_SIGNING_UP.String(),
				NumberOfTries: 0,
				Payload:       map[string]any{enum.QueueFields_EMAIL.String(): "test1@gmail.com"},
				ServiceName:   m.DbName,
			}

			dataChan <- &[]types.Data{authData}
			return

		case enum.ServiceName_TRACKER_SERVICE.String():

			payload := map[string]any{enum.QueueFields_USER_ID.String(): "123", enum.QueueFields_PLAN_NAME.String(): "abc", enum.QueueFields_EXERCISE_NAMES.String(): []string{"exer_1", "exer_2"}}

			trackerData := types.Data{
				DbId:          "123",
				TargetService: enum.ServiceName_PLAN_SERVICE.String(),
				CreatedBy:     enum.ServiceName_TRACKER_SERVICE.String(),
				Task:          enum.TaskName_UPDATE_PLAN.String(),
				NumberOfTries: 0,
				Payload:       payload,
				ServiceName:   m.DbName,
			}

			dataChan <- &[]types.Data{trackerData}
			return
		}
	}

	if !m.HasData {

		data := types.Data{
			NoData:      true,
			ServiceName: m.DbName,
		}

		dataChan <- &[]types.Data{data}
		return
	}

}

func (m *MockDb) UpdateTaskStatus(ctx context.Context, dbIndex string, updateValue string) error {
	if m.Down {
		return fmt.Errorf("db is down")
	}

	m.Update[0] = updateValue

	return nil
}

func (m *MockDb) UpdateTaskStatusWithNumberOfTries(ctx context.Context, dbIndex string, updateValue string) error {
	if m.Down {
		return fmt.Errorf("db is down")
	}

	return nil
}
