package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type JsonArray []interface{}

func (obj *JsonArray) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		return json.Unmarshal(v, &obj)
	case string:
		return json.Unmarshal([]byte(v), &obj)
	case []interface{}:
		*obj = v
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (obj JsonArray) Value() (driver.Value, error) {
	return json.Marshal(obj)
}

func (obj JsonArray) Unmarshal(to interface{}) error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, to)
}
