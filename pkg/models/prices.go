package models

import (
	"time"
)

type HistoricPrice struct {
	Id          int64     `db:"id" json:"-"`
	BoardgameId int64     `db:"boardgame_id" json:"boardgame_id"`
	CrDate      time.Time `db:"cr_date" json:"cr_date"`
	Price       float64   `db:"price" json:"price"`
	Stock       int       `db:"stock" json:"stock"`
	StoreId     int64     `db:"store_id" json:"store_id"`
}

type Price struct {
	Id             int64           `db:"id" json:"id"`
	CrDate         time.Time       `db:"cr_date" json:"cr_date"`
	Name           string          `db:"name" json:"name"`
	BoardgameId    JsonNullInt64   `db:"boardgame_id" json:"boardgame_id"`
	BoardgameName  string          `db:"boardgame_name" json:"boardgame_name"`
	BoardgameThumb string          `db:"thumb" json:"thumb"`
	StoreId        int64           `db:"store_id" json:"store_id"`
	StoreName      string          `db:"store_name" json:"-"`
	StoreThumb     string          `db:"store_thumb" json:"store_thumb"`
	Price          float64         `db:"price" json:"price"`
	Stock          int             `db:"stock" json:"stock"`
	Url            string          `db:"url" json:"url"`
	Levenshtein    int             `db:"levenshtein" json:"-"`
	Hamming        int             `db:"hamming" json:"-"`
	Rank           JsonNullInt64   `db:"rank" json:"rank"`
	Batch          int64           `db:"batch" json:"-"`
	Mapped         bool            `db:"mapped" json:"-"`
	HistoricPrices []HistoricPrice `json:"-"`
}
