package types

import (
	"fmt"
	"time"
)

type Data struct {
	DbId          string         `db:"id"`
	TargetService string         `db:"target_service"`
	Task          string         `db:"task"`
	Status        string         `db:"status"`
	Payload       map[string]any `db:"payload"`
	CreatedAt     time.Time      `db:"created_at"`
	NumberOfTries int            `db:"number_of_tries"`
}

func (d *Data) GetEmail() (string, error) {
	email, ok := d.Payload["email"].(string)
	if !ok {
		return "", fmt.Errorf("error getting email")
	}

	return email, nil

}