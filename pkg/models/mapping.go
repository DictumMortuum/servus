package models

type Mapping struct {
	Id          int64  `db:"id" json:"id"`
	BoardgameId int64  `db:"boardgame_id" json:"boardgame_id"`
	Name        string `db:"name" json:"name"`
}
