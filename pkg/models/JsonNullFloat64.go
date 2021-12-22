package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type JsonNullFloat64 sql.NullFloat64

func NewJsonNullFloat64(i float64) JsonNullFloat64 {
	return JsonNullFloat64{
		Float64: i,
		Valid:   true,
	}
}

func (obj *JsonNullFloat64) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &obj)
		return nil
	case nil:
		json.Unmarshal(nil, &obj)
		return nil
	case float64:
		obj.Valid = true
		obj.Float64 = v
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (obj JsonNullFloat64) Value() (driver.Value, error) {
	if !obj.Valid {
		return nil, nil
	}

	return obj.Float64, nil
}

func (v JsonNullFloat64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullFloat64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Float64 = *x
	} else {
		v.Valid = false
	}
	return nil
}
