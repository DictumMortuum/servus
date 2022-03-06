package models

import (
	"time"
)

type Price struct {
	Id             int64         `db:"id" json:"id"`
	CrDate         time.Time     `db:"cr_date" json:"cr_date"`
	Name           string        `db:"name" json:"name"`
	BoardgameId    JsonNullInt64 `db:"boardgame_id" json:"boardgame_id"`
	BoardgameName  string        `db:"boardgame_name" json:"boardgame_name"`
	BoardgameThumb string        `db:"thumb" json:"thumb"`
	StoreId        int64         `db:"store_id" json:"-"`
	StoreName      string        `db:"store_name" json:"store_name"`
	StoreThumb     string        `db:"store_thumb" json:"store_thumb"`
	Price          float64       `db:"price" json:"price"`
	Stock          bool          `db:"stock" json:"stock"`
	Url            string        `db:"url" json:"url"`
	Levenshtein    int           `db:"levenshtein" json:"levenshtein"`
	Hamming        int           `db:"hamming" json:"hamming"`
	Rank           JsonNullInt64 `db:"rank" json:"rank"`
	Batch          int64         `db:"batch" json:"-"`
	Mapped         bool          `db:"mapped" json:"-"`
}
