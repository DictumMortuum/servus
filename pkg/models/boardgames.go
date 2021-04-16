package models

import (
	"time"
)

type Boardgame struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type Player struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type Play struct {
	Id          int64     `db:"id" json:"id"`
	CrDate      time.Time `db:"cr_date" json:"cr_date"`
	Date        time.Time `db:"date" json:"date"`
	BoardgameId int64     `db:"boardgame_id"`
}

type WingspanPlay struct {
	Id               int64     `db:"id" json:"id"`
	PlayId           int64     `db:"play_id" json:"play_id"`
	PlayerId         int64     `db:"player_id" json:"player_id"`
	Date             time.Time `json:"date"`
	PlayerName       string    `db:"player_name" json:"player"`
	BirdPoints       int       `db:"bird_points" json:"birds"`
	BonusPoints      int       `db:"bonus_points" json:"bonus"`
	EndofroundPoints int       `db:"endofround_points" json:"endofround"`
	EggPoints        int       `db:"egg_points" json:"egg"`
	FoodPoints       int       `db:"food_points" json:"food"`
	TuckedPoints     int       `db:"tucked_points" json:"tucked"`
}

type DuelPlay struct {
	Id             int64  `db:"id" json:"id"`
	PlayId         int64  `db:"play_id" json:"play_id"`
	PlayerId       int64  `db:"player_id" json:"player_id"`
	PlayerName     string `db:"player_name" json:"player"`
	Expansions     string `db:"expansions" json:"expansions"`
	BluePoints     int    `db:"blue_points" json:"blue"`
	GreenPoints    int    `db:"green_points" json:"green"`
	YellowPoints   int    `db:"yellow_points" json:"yellow"`
	PurplePoints   int    `db:"purple_points" json:"purple"`
	WonderPoints   int    `db:"wonder_points" json:"wonder"`
	MarkerPoints   int    `db:"marker_points" json:"marker"`
	CoinPoints     int    `db:"coin_points" json:"coin"`
	BattlePoints   int    `db:"battle_points" json:"battle"`
	PantheonPoints int    `db:"pantheon_points" json:"pantheon"`
	BattleVictory  bool   `db:"battle_victory" json:"battle_victory"`
	ScienceVictory bool   `db:"science_victory" json:"science_victory"`
}
