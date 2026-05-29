package utils

import (
	"encoding/json"
	"fmt"

	"github.com/sakamoto-max/rabbit_mq/types"
)

func ConvertToJson(src []byte) (map[string]any, error) {

	var data map[string]any

	err := json.Unmarshal(src, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to map[string]string")
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
