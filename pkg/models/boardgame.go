package models

import (
	"database/sql"
)

type Boardgame struct {
	Id       int64          `db:"id" json:"id"`
	Name     string         `db:"name" json:"name"`
	Unmapped bool           `db:"unmapped" json:"unmapped"`
	AtlasId  sql.NullString `db:"atlas_id" json:"atlas_id"`
	Atlas    []AtlasResult  `db:"atlas" json:"atlas"`
}

func (obj Boardgame) Insert() string {
	return `insert into tboardgames (name,unmapped) values (:name,0)`
}

func (obj Boardgame) Select() string {
	return `select * from tboardgames where name = :name and unmapped != 1`
}

func (obj Boardgame) Exists() string {
	return `select id from tboardgames where name = :name`
}

// type WingspanPlay struct {
// 	Id               int64     `db:"id" json:"id"`
// 	PlayId           int64     `db:"play_id" json:"play_id"`
// 	PlayerId         int64     `db:"player_id" json:"player_id"`
// 	Data             Json      `db:"data"`
// 	Date             time.Time `json:"date"`
// 	PlayerName       string    `db:"player_name" json:"player"`
// 	BirdPoints       int       `db:"bird_points" json:"birds"`
// 	BonusPoints      int       `db:"bonus_points" json:"bonus"`
// 	EndofroundPoints int       `db:"endofround_points" json:"endofround"`
// 	EggPoints        int       `db:"egg_points" json:"egg"`
// 	FoodPoints       int       `db:"food_points" json:"food"`
// 	TuckedPoints     int       `db:"tucked_points" json:"tucked"`
// }
