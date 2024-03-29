package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/DictumMortuum/servus/pkg/w3m"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func ScrapeGamesCom(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(18)
	detected := 0

	conn, ch, q, err := rabbitmq.SetupQueue("prices")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer ch.Close()

	rows, err := updateBatch(db, store_id)
	if err != nil {
		return nil, err
	}

	log.Printf("Scraper %d started - resetting %d rows\n", store_id, rows)

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	c.OnHTML("a.ty-pagination__next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		local_link, _ := w3m.Download(link)
		c.Visit(local_link)
	})

	c.OnHTML(".col-tile", func(e *colly.HTMLElement) {
		var stock int
		raw_price := e.ChildText(".ty-price")

		if childHasClass(e, "button.ty-btn__primary", "ty-btn__add-to-cart") {
			stock = 0
		} else {
			stock = 2
		}

		item := models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".cm-image", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".product-title", "href"),
		}

		detected++
		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	for _, url := range []string{"https://www.gamescom.gr/epitrapezia-el/", "https://www.gamescom.gr/epitrapezia-el/category-124/"} {
		local, err := w3m.Download(url)
		if err != nil {
			return nil, err
		}

		c.Visit(local)
	}

	c.Wait()

	return map[string]interface{}{
		"name":     "Gamescom",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
