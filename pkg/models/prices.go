package models

import (
	"time"
)

type Price struct {
	Id            int64     `db:"id" json:"id"`
	CrDate        time.Time `db:"cr_date" json:"cr_date"`
	Name          string    `db:"name" json:"name"`
	BoardgameId   int64     `db:"boardgame_id" json:"boardgame_id"`
	BoardgameName string    `db:"boardgame_name" json:"boardgame_name"`
	StoreId       int64     `db:"store_id" json:"store_id"`
	StoreName     string    `db:"store_name" json:"store_name"`
	Price         float64   `db:"price" json:"price"`
	Stock         bool      `db:"stock" json:"stock"`
	Url           string    `db:"url" json:"url"`
	Distance      int       `db:"distance" json:"distance"`
	Rank          int64     `db:"rank" json:"rank"`
}
