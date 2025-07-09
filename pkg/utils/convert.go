package utils

import (
	"encoding/json"
	"errors"
)

func StructToMap(data interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

	if data == nil {
		return nil, errors.New("input data is nil")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		result = make(map[string]interface{})
	}
	if len(result) == 0 {
		return nil, errors.New("converted map is empty")
	}

	return result, nil
}
