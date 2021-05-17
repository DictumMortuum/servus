package models

import (
	"database/sql"
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

type JsonNullInt64 sql.NullInt64

func (obj *JsonNullInt64) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &obj)
		return nil
	case nil:
		json.Unmarshal(nil, &obj)
		return nil
	case int64:
		obj.Valid = true
		obj.Int64 = v
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (obj JsonNullInt64) Value() (driver.Value, error) {
	if !obj.Valid {
		return nil, nil
	}

	return obj.Int64, nil
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}
