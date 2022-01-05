package models

import (
	"time"
)

type Boardgame struct {
	Id    int64           `db:"id" json:"id"`
	Name  string          `db:"name" json:"name"`
	Date  time.Time       `json:"validUntil"`
	Data  Json            `db:"data" json:"data"`
	Guid  JsonNullString  `db:"tx_guid" json:"tx_guid"`
	Cost  JsonNullFloat64 `db:"cost" json:"cost"`
	Rank  JsonNullInt64   `db:"rank" json:"rank"`
	Thumb JsonNullString  `db:"thumb" json:"thumb"`
}
