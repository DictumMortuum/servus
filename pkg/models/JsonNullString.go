package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type JsonNullString sql.NullString

func NewJsonNullString(i string) JsonNullString {
	return JsonNullString{
		String: i,
		Valid:  true,
	}
}

func (obj *JsonNullString) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		obj.Valid = true
		obj.String = string(v)
		return nil
	case nil:
		obj.Valid = false
		obj.String = ""
		return nil
	case string:
		obj.Valid = true
		obj.String = v
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (obj JsonNullString) Value() (driver.Value, error) {
	if !obj.Valid {
		return nil, nil
	}

	return obj.String, nil
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullString) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}
	return nil
}
