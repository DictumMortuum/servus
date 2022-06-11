package search

import (
	"encoding/xml"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
)

type madnessRoot struct {
	XMLName  xml.Name         `xml:"products"`
	Products []madnessProduct `xml:"product"`
}

type madnessProduct struct {
	XMLName  xml.Name `xml:"product"`
	SKU      string   `xml:"id"`
	Name     string   `xml:"title"`
	ThumbUrl string   `xml:"image_link"`
	Price    string   `xml:"price"`
	Stock    string   `xml:"availability"`
	Link     string   `xml:"link"`
}

func madnessAvailbilityToStock(s string) int {
	switch s {
	case "in stock":
		return 0
	case "on backorder":
		return 1
	case "out of stock":
		return 2
	default:
		return 2
	}
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

	rs := madnessRoot{}
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
			Stock:      madnessAvailbilityToStock(item.Stock),
			Price:      getPrice(item.Price),
			Url:        item.Link,
		})
	}

	err = updateBatch(db, 16)
	if err != nil {
		return nil, err
	}

	for _, item := range prices {
		err = UpsertPrice(db, item)
		if err != nil {
			return nil, err
		}
	}

	return prices, nil
}
