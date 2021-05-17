package models

import (
	"time"
)

type Scrapable interface {
	Scrape() ([]ScraperData, error)
	GetStore() Store
	ScrapePrices() ([]ScraperPrice, error)
}

type ScraperData struct {
	Id          int64         `db:"id" json:"id"`
	StoreId     JsonNullInt64 `db:"store_id" json:"store_id"`
	Store       string        `json:"store"`
	BoardgameId JsonNullInt64 `db:"boardgame_id" json:"boardgame_id"`
	CrDate      time.Time     `db:"cr_date" json:"cr_date"`
	Title       string        `db:"title" json:"title"`
	Link        string        `db:"link" json:"link"`
	SKU         string        `db:"sku" json:"sku"`
	Active      time.Time     `db:"active" json:"active"`
}

func (obj ScraperData) Insert() string {
	return `
	insert into tboardgamescraperdata (
		store_id,
		boardgame_id,
		cr_date,
		title,
		link,
		sku,
		active
	) values (
		:store_id,
		NULL,
		NOW(),
		:title,
		:link,
		:sku,
		NOW()
	)`
}

func (obj ScraperData) Select() string {
	return `select * from tboardgamescraperdata where store_id = :store_id and title = :title and sku = :sku`
}

func (obj ScraperData) Exists() string {
	return `select 1 from tboardgamescraperdata where store_id = :store_id and title = :title and sku = :sku`
}

type ScraperPrice struct {
	Id      int64     `db:"id" json:"id"`
	DataId  int64     `db:"data_id" json:"data_id"`
	CrDate  time.Time `db:"cr_date" json:"cr_date"`
	Price   float64   `db:"price" json:"price"`
	InStock bool      `db:"in_stock" json:"in_stock"`
}

func (obj ScraperPrice) Insert() string {
	return `insert into tboardgamescraperprices (data_id, cr_date, price, in_stock) values (:data_id, NOW(), :price, :in_stock)`
}

func (obj ScraperPrice) Select() string {
	return `select * from tboardgamescraperprices where data_id = :data_id and price = :price and in_stock = :in_stock`
}

func (obj ScraperPrice) Exists() string {
	return `select 1 from tboardgamescraperprices where data_id = :data_id and price = :price and in_stock = :in_stock`
}

type ScraperResult struct {
	Id          int64         `db:"id" json:"id"`
	CrDate      time.Time     `db:"cr_date" json:"cr_date"`
	Title       string        `db:"title" json:"title"`
	Store       string        `db:"store" json:"store"`
	BoardgameId JsonNullInt64 `db:"boardgame_id" json:"boardgame_id"`
	Link        string        `db:"link" json:"link"`
	SKU         string        `db:"sku" json:"sku"`
	Price       float64       `db:"price" json:"price"`
	InStock     bool          `db:"in_stock" json:"in_stock"`
}
