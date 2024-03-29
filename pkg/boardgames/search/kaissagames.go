package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeKaissaGames(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(9)
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

	collector := colly.NewCollector(
		colly.AllowedDomains("kaissagames.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML("li.item.product-item", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price")

		var stock int

		if e.ChildText(".release-date") != "" {
			stock = 1
		} else {
			if !childHasClass(e, "div.stock", "unavailable") {
				stock = 0
			} else {
				stock = 2
			}
		}

		item := models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".product-image-photo", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		}

		detected++
		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://kaissagames.com/b2c_gr/xenoglossa-epitrapezia.html")
	collector.Wait()

	return map[string]interface{}{
		"name":     "Kaissa Games",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
