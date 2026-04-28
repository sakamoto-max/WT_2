package utils

import (
	"encoding/json"
	"fmt"

	"github.com/sakamoto-max/rabbit_mq/types"
)

func GetUserId(d *types.Data) (string, error) {
	str, ok := d.Payload["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("error getting user id")
	}

	return str, nil
}

func GetPlanName(d *types.Data) (string, error) {
	str, ok := d.Payload["plan_name"].(string)
	if !ok {
		return "", fmt.Errorf("error getting planName")
	}

	return str, nil
}

func GetNewExercises(d *types.Data) ([]string, error) {
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

func GetEmail(d *types.Data) (string, error) {
	data, ok := d.Payload["email"].(string)
	if !ok {
		return "", fmt.Errorf("error getting the email")
	}

	return data, nil
}

func ConvertIntoBytes(payload any) (*[]byte, error) {

	dataInBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error in converting data into bytes : %w", err)
	}

	return &dataInBytes, nil
}

func ConvertIntoStruct(data *[]byte) *types.Data {

	var D types.Data

	_ = json.Unmarshal(*data, &D)

	return &D
}
