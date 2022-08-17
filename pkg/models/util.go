package models

import (
	"fmt"
	"time"
)

func getTime(data map[string]interface{}, key string) (time.Time, error) {
	if val, ok := data["date"]; ok {
		//"Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"
		t, err := time.Parse("2006-01-02", val.(string))
		if err != nil {
			return time.Now(), err
		}

		return t, nil
	} else {
		return time.Now(), fmt.Errorf("please provide a '%s' parameter", key)
	}
}

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

func getFloat(data map[string]interface{}, key string) (float64, error) {
	if val, ok := data[key]; ok {
		return val.(float64), nil
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
