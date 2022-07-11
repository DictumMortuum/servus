package search

import (
	"encoding/xml"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
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
	store_id := int64(16)
	detected := 0

	rconn, ch, q, err := rabbitmq.SetupQueue("prices")
	if err != nil {
		return nil, err
	}
	defer rconn.Close()
	defer ch.Close()

	rows, err := updateBatch(db, store_id)
	if err != nil {
		return nil, err
	}

	log.Printf("Scraper %d started - resetting %d rows\n", store_id, rows)

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

	for _, item := range rs.Products {
		item := models.Price{
			Name:       item.Name,
			StoreId:    store_id,
			StoreThumb: item.ThumbUrl,
			Stock:      madnessAvailbilityToStock(item.Stock),
			Price:      getPrice(item.Price),
			Url:        item.Link,
		}

		detected++
		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"name":     "Boards of Madness",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
