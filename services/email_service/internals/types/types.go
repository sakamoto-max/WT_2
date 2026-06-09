package types

import (
	"fmt"

	"github.com/sakamoto-max/wt_2_proto/shared/enum"
)

type Data struct {
	DbId          string         `json:"dbId"`
	TargetService string         `json:"targetService"`
	TaskName      string         `json:"taskName"`
	Status        string         `json:"status"`
	Payload       map[string]any `json:"payload"`
	SentBy        string         `json:"sentBy"`
	Err           error          `json:"err"`
}

func (d *Data) GetEmail() (string, error) {
	email, ok := d.Payload[enum.QueueFields_EMAIL.String()].(string)
	if !ok {
		return "", fmt.Errorf("error getting email")
	}

	return email, nil
}

func (d Data) Failed(err error) Data {
	return Data{
		DbId:          d.DbId,
		TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
		Status:        enum.TaskStatus_TASK_FAILED.String(),
		SentBy:        d.TargetService,
		TargetService: d.SentBy,
		Err:           err,
	}
}

func (d Data) Succeded() Data {
	return Data{
		DbId:          d.DbId,
		TaskName:      enum.TaskName_UPDATE_VALUE_IN_DB.String(),
		Status:        enum.TaskStatus_TASK_COMPLETED.String(),
		SentBy:        d.TargetService,
		TargetService: d.SentBy,
	}
}
