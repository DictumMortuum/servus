package models

import (
	"database/sql"
)

type Message struct {
	Id       int64        `db:"id"`
	Msg      string       `db:"msg"`
	DateSend sql.NullTime `db:"date_send"`
}
