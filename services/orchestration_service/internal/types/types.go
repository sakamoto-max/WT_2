package types

import (
	"fmt"
	"time"
	"github.com/sakamoto-max/wt_2-pkg/utils"
)

type Data struct {
	DbId          string    `db:"id"`
	TargetService string    `db:"target_service"`
	CreatedBy     string    `db:"created_by"`
	Task          string    `db:"task"`
	Status        string    `db:"status"`
	Payload       any       `db:"payload"`
	CreatedAt     time.Time `db:"created_at"`
	NumberOfTries int       `db:"number_of_tries"`
}

func (d *Data) ConvertToBytes() (*[]byte, error) {

	dataInBytes, err := utils.ConvertIntoBytes(d)
	if err != nil {
		return nil, fmt.Errorf("error while converting data into bytes : %w", err)
	}

	return dataInBytes, nil
}

type taskStatus string

var (
	TaskCompleted    taskStatus = "completed"
	TaskPending      taskStatus = "pending"
	TaskNotCompleted taskStatus = "not_completed"
	TaskFailed       taskStatus = "failed"
)

type serviceName string

var (
	PlanService     serviceName = "plan_service"
	AuthService     serviceName = "auth_service"
	TrackerService  serviceName = "tracker_service"
	ExerciseService serviceName = "exercise_service"
	EmailService    serviceName = "email_service"
)
