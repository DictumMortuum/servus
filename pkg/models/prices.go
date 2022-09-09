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
	Id               int64           `db:"id" json:"id"`
	Key              string          `json:"key"`
	CrDate           time.Time       `db:"cr_date" json:"cr_date"`
	Name             string          `db:"name" json:"name"`
	TransformedName  string          `json:"transformed_name"`
	BoardgameId      JsonNullInt64   `db:"boardgame_id" json:"boardgame_id"`
	BoardgameName    string          `db:"boardgame_name" json:"boardgame_name"`
	BoardgameThumb   string          `db:"thumb" json:"thumb"`
	BoardgamePreview string          `db:"preview" json:"preview"`
	StoreId          int64           `db:"store_id" json:"store_id"`
	StoreName        string          `db:"store_name" json:"-"`
	StoreThumb       string          `db:"store_thumb" json:"store_thumb"`
	Price            float64         `db:"price" json:"price"`
	Stock            int             `db:"stock" json:"stock"`
	Url              string          `db:"url" json:"url"`
	Levenshtein      int             `db:"levenshtein" json:"-"`
	ExtraId          int             `db:"extra_id" json:"-"`
	Rank             JsonNullInt64   `db:"rank" json:"rank"`
	Batch            int64           `db:"batch" json:"-"`
	Mapped           bool            `db:"mapped" json:"-"`
	Ignored          bool            `db:"ignored"`
	HistoricPrices   []HistoricPrice `json:"-"`
	ProductId        string          `json:"-"`
}

func (rs *Price) Constructor() []func(map[string]interface{}, bool) error {
	return []func(map[string]interface{}, bool) error{
		rs.SetStoreId,
		rs.SetBoardgameId,
		rs.SetPrice,
		rs.SetStock,
		rs.SetUrl,
		rs.SetStoreThumb,
		rs.SetBatch,
		rs.SetLevenshtein,
		rs.SetMapped,
		rs.SetIgnored,
	}
}

func (rs *Price) SetBoardgameId(data map[string]interface{}, create bool) error {
	id, err := getJsonNullInt64(data, "boardgame_id")
	if err != nil {
		return err
	}

	rs.BoardgameId = id
	return nil
}

func (rs *Price) SetStoreId(data map[string]interface{}, create bool) error {
	id, err := getInt64(data, "store_id")
	if err != nil {
		return err
	}

	rs.StoreId = id
	return nil
}

func (rs *Price) SetPrice(data map[string]interface{}, create bool) error {
	price, err := getFloat(data, "price")
	if err != nil {
		return err
	}

	rs.Price = price
	return nil
}

func (rs *Price) SetStock(data map[string]interface{}, create bool) error {
	stock, err := getInt(data, "stock")
	if err != nil {
		return err
	}

	rs.Stock = stock
	return nil
}

func (rs *Price) SetUrl(data map[string]interface{}, create bool) error {
	url, err := getString(data, "url")
	if err != nil {
		return err
	}

	rs.Url = url
	return nil
}

func (rs *Price) SetStoreThumb(data map[string]interface{}, create bool) error {
	store_thumb, err := getString(data, "store_thumb")
	if err != nil {
		return err
	}

	rs.StoreThumb = store_thumb
	return nil
}

func (rs *Price) SetBatch(data map[string]interface{}, create bool) error {
	batch, err := getInt64(data, "batch")
	if err != nil {
		return err
	}

	rs.Batch = batch
	return nil
}

func (rs *Price) SetLevenshtein(data map[string]interface{}, create bool) error {
	rs.Levenshtein = 0
	return nil
}

func (rs *Price) SetMapped(data map[string]interface{}, create bool) error {
	mapped, err := getBool(data, "mapped")
	if err != nil {
		return err
	}

	rs.Mapped = mapped
	return nil
}

func (rs *Price) SetIgnored(data map[string]interface{}, create bool) error {
	ignored, err := getBool(data, "ignored")
	if err != nil {
		return err
	}

	rs.Ignored = ignored
	return nil
}
