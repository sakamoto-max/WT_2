package types

import (
	"fmt"

	"github.com/sakamoto-max/wt_2_proto/shared/enum"
)

type Data struct {
	DbId          string `db:"id"`
	SentBy        string
	TaskName      string         `db:"task"`
	Payload       map[string]any `db:"payload"`
	TargetService string
}

func (d *Data) GetUserId() (string, error) {
	str, ok := d.Payload[enum.QueueFields_USER_ID.String()].(string)
	if !ok {
		return "", fmt.Errorf("error getting user id")
	}

	return str, nil
}

func (d *Data) GetPlanName() (string, error) {
	str, ok := d.Payload[enum.QueueFields_PLAN_NAME.String()].(string)
	if !ok {
		return "", fmt.Errorf("error getting planName")
	}

	return str, nil

}
func (d *Data) GetNewExercises() ([]string, error) {
	raw, ok := d.Payload[enum.QueueFields_EXERCISE_NAMES.String()].([]any)
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
