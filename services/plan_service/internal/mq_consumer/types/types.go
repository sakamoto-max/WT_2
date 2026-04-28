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

func (d *Data) GetUserId() (string, error) {
	str, ok := d.Payload["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("error getting user id")
	}

	return str, nil
}
func (d *Data) GetPlanName() (string, error) {
	str, ok := d.Payload["plan_name"].(string)
	if !ok {
		return "", fmt.Errorf("error getting planName")
	}

	return str, nil

}
func (d *Data) GetNewExercises() ([]string, error) {
	raw, ok := d.Payload["exercise_names"].([]any)
	if !ok {
		return nil, fmt.Errorf("error getting exercise_names")
	}

	var data []string

	for _, v := range raw {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("one of the element is not string")
		}

		data = append(data, str)
	}

	return data, nil
}
func (d *Data) GetEmail() (string, error) {
	data, ok := d.Payload["email"].(string)
	if !ok {
		return "", fmt.Errorf("error getting the email")
	}

	return data, nil
}
