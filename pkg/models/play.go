package models

import (
	"time"
)

type Play struct {
	Id          int64     `db:"id" json:"id"`
	BoardgameId int64     `db:"boardgame_id" json:"boardgame_id"`
	CrDate      time.Time `db:"cr_date" json:"cr_date"`
	Date        time.Time `db:"date" json:"date"`
	Boardgame   string    `db:"name" json:"boardgame"`
	Stats       []Stats   `json:"stats"`
	Probability float64   `json:"probability"`
}
