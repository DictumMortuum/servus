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

type hobbyRoot struct {
	XMLName xml.Name      `xml:"mywebstore"`
	Store   hobbyProducts `xml:"products"`
}

type hobbyProducts struct {
	XMLName  xml.Name  `xml:"products"`
	Products []product `xml:"product"`
}

type product struct {
	XMLName      xml.Name `xml:"product"`
	SKU          string   `xml:"id"`
	Name         string   `xml:"name"`
	ThumbUrl     string   `xml:"image"`
	Category     string   `xml:"category"`
	Price        string   `xml:"price_with_vat"`
	Stock        string   `xml:"instock"`
	Availability string   `xml:"availability"`
	Link         string   `xml:"link"`
}

func ScrapeHobbyTheory(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(23)

	log.Printf("Scraper %d started\n", store_id)

	rconn, ch, q, err := rabbitmq.SetupQueue("prices")
	if err != nil {
		return nil, err
	}
	defer rconn.Close()
	defer ch.Close()

	err = updateBatch(db, store_id)
	if err != nil {
		return nil, err
	}

	link := "https://feed.syntogether.com/skroutz/xml?shop=hobbytheory.myshopify.com"
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

	rs := hobbyRoot{}
	err = xml.Unmarshal(body, &rs)
	if err != nil {
		return nil, err
	}

	categories := []string{
		"Επιτραπέζια Παιχνίδια Οικογενειακά",
		"Επιτραπέζια Παιχνίδια Παρέας",
		"Επιτραπέζια Παιχνίδια Πολέμου",
		"Επιτραπέζια Παιχνίδια Στρατηγικής",
		"Θεματικά Επιτραπέζια Παιχνίδια",
	}

	for _, item := range rs.Store.Products {
		for _, category := range categories {
			if item.Category == category {
				var stock int

				if item.Stock == "Y" {
					stock = 0
				} else {
					stock = 2
				}

				item := models.Price{
					Name:       item.Name,
					StoreId:    store_id,
					StoreThumb: item.ThumbUrl,
					Stock:      stock,
					Price:      getPrice(item.Price),
					Url:        item.Link,
				}

				err = rabbitmq.InsertQueueItem(ch, q, item)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	return nil, nil
}
