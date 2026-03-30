package types

import "time"

type Data struct {
	Id            int            `db:"id"`
	TargetService string         `db:"target_service"`
	Task          string         `db:"task"`
	Status        string         `db:"status"`
	Payload       map[string]int `db:"payload"`
	CreatedAt     time.Time      `db:"created_at"`
	NumberOfTries *int           `db:"number_of_tries"`
}

type queueName string

var (
	PlanQueue queueName = "plan_queue"
)
