package models

import (
	"time"
)

type KeyVal struct {
	Id   string    `db:"id"`
	Date time.Time `db:"date"`
	Data Json      `db:"json"`
}

func (obj KeyVal) Unmarshal(to interface{}) error {
	return obj.Data.Unmarshal(to)
}
