package models

import (
	"fmt"
)

func getInt64(data map[string]interface{}, key string) (int64, error) {
	if val, ok := data[key]; ok {
		return int64(val.(float64)), nil
	} else {
		return -1, fmt.Errorf("please provide a '%s' parameter", key)
	}
}

func getInt(data map[string]interface{}, key string) (int, error) {
	if val, ok := data[key]; ok {
		return int(val.(float64)), nil
	} else {
		return -1, fmt.Errorf("please provide a '%s' parameter", key)
	}
}

func getString(data map[string]interface{}, key string) (string, error) {
	if val, ok := data[key]; ok {
		return val.(string), nil
	} else {
		return "", fmt.Errorf("please provide a '%s' parameter", key)
	}
}

func getBool(data map[string]interface{}, key string) (bool, error) {
	if val, ok := data[key]; ok {
		return val.(bool), nil
	} else {
		return false, fmt.Errorf("please provide a '%s' parameter", key)
	}
}

func getJsonNullInt64(data map[string]interface{}, key string) (JsonNullInt64, error) {
	if val, ok := data[key]; ok {
		return JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}, nil
	} else {
		return JsonNullInt64{
			Int64: -1,
			Valid: false,
		}, fmt.Errorf("please provide a '%s' parameter", key)
	}
}

func getJsonNullString(data map[string]interface{}, key string) (JsonNullString, error) {
	if val, ok := data[key]; ok {
		return JsonNullString{
			String: val.(string),
			Valid:  true,
		}, nil
	} else {
		return JsonNullString{
			String: "",
			Valid:  false,
		}, fmt.Errorf("please provide a '%s' parameter", key)
	}
}
