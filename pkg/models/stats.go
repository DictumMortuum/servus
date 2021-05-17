package models

type Stats struct {
	Id       int64 `db:"id" json:"id"`
	PlayId   int64 `db:"play_id" json:"play_id"`
	PlayerId int64 `db:"player_id" json:"player_id"`
	Data     Json  `db:"data" json:"data"`
}
