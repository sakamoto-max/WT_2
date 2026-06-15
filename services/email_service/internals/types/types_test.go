package types

import (
	"fmt"
	"testing"

	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/stretchr/testify/assert"
)

func Test_GetEmail(t *testing.T) {

	tests := []struct {
		name    string
		email   any
		wantErr bool
	}{
		{name: "success", email: "test2@gmail.com"},
		{name: "fails", email: 123456789, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := Data{
				Payload: map[string]any{
					enum.QueueFields_EMAIL.String(): test.email,
				},
			}
			emailGot, err := data.GetEmail()
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, emailGot, test.email)

		})
	}

}

func Test_Failed(t *testing.T) {

	err := fmt.Errorf("some error")

	data := Data{
		TaskName:      "some task",
		Status:        "some status",
		SentBy:        "sent",
		TargetService: "target",
	}

	data2 := data.Failed(err)

	assert.Error(t, data2.Err)
	assert.Equal(t, data2.TaskName, enum.TaskName_UPDATE_VALUE_IN_DB.String())
	assert.Equal(t, data2.SentBy, data.TargetService)
	assert.Equal(t, data2.Status, enum.TaskStatus_TASK_FAILED.String())
	assert.Equal(t, data2.TargetService, data.SentBy)
}

func Test_Succeded(t *testing.T) {

	data := Data{
		TaskName:      "some task",
		Status:        "some status",
		SentBy:        "sent",
		TargetService: "target",
	}

	data2 := data.Succeded()

	assert.NoError(t, data2.Err)
	assert.Equal(t, data2.TaskName, enum.TaskName_UPDATE_VALUE_IN_DB.String())
	assert.Equal(t, data2.SentBy, data.TargetService)
	assert.Equal(t, data2.TargetService, data.SentBy)
	assert.Equal(t, data2.Status, enum.TaskStatus_TASK_COMPLETED.String())

}
