package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Json map[string]interface{}

func (obj *Json) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		return json.Unmarshal(v, &obj)
	case string:
		return json.Unmarshal([]byte(v), &obj)
	case map[string]interface{}:
		*obj = v
		return nil
	case nil:
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (obj Json) Value() (driver.Value, error) {
	return json.Marshal(obj)
}

func (obj Json) Unmarshal(to interface{}) error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, to)
}
