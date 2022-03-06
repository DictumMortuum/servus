package search

import (
	"encoding/xml"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
)

type Root struct {
	XMLName  xml.Name  `xml:"products"`
	Products []Product `xml:"product"`
}

type Product struct {
	XMLName  xml.Name `xml:"product"`
	SKU      string   `xml:"id"`
	Name     string   `xml:"title"`
	ThumbUrl string   `xml:"image_link"`
	Price    string   `xml:"price"`
	Stock    string   `xml:"availability"`
	Link     string   `xml:"link"`
}

func ScrapeBoardsOfMadness(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	link := "https://boardsofmadness.com/wp-content/uploads/woo-product-feed-pro/xml/sVVFMsJLyEEtvbil4fbIOdm8b4ha7ewz.xml"
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	conn := &http.Client{}
	resp, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rs := Root{}
	err = xml.Unmarshal(body, &rs)
	if err != nil {
		return nil, err
	}

	prices := []models.Price{}
	for _, item := range rs.Products {
		prices = append(prices, models.Price{
			Name:       item.Name,
			StoreId:    16,
			StoreThumb: item.ThumbUrl,
			Stock:      item.Stock == "in stock",
			Price:      getPrice(item.Price),
			Url:        item.Link,
		})
	}

	for _, item := range prices {
		item.BoardgameId = models.JsonNullInt64{
			Int64: -1,
			Valid: false,
		}

		exists, err := exists(db, item)
		if err != nil {
			return nil, err
		}

		if !exists {
			_, err := create(db, item)
			if err != nil {
				return nil, err
			}
		}
	}

	return prices, nil
}
