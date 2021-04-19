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
		json.Unmarshal(v, &obj)
		return nil
	case string:
		json.Unmarshal([]byte(v), &obj)
		return nil
	case nil:
		json.Unmarshal(nil, &obj)
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (obj Json) Value() (driver.Value, error) {
	return json.Marshal(obj)
}
