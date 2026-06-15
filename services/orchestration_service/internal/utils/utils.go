package utils

import (
	"encoding/json"
	"fmt"
)

func ConvertToJson(src []byte) (map[string]any, error) {

	var data map[string]any

	err := json.Unmarshal(src, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to map[string]string")
	}

	fmt.Println("data", data)

	return data, nil
}
