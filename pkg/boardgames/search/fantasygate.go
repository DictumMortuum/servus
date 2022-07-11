package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeFantasyGate(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(2)
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
		colly.AllowedDomains("www.fantasygate.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".sblock4", func(e *colly.HTMLElement) {
		var stock int
		raw_price := e.ChildText(".jshop_price")

		if childHasClass(e, ".btn", "button_buy") {
			stock = 0
		} else {
			stock = 2
		}

		item := models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".jshop_img", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
		}

		detected++
		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.Post("https://www.fantasygate.gr/strategygames", map[string]string{
		"limit": "99999",
	})

	collector.Post("https://www.fantasygate.gr/family-games", map[string]string{
		"limit": "99999",
	})

	collector.Post("https://www.fantasygate.gr/cardgames", map[string]string{
		"limit": "99999",
	})

	collector.Wait()

	return map[string]interface{}{
		"name":     "Fantasy Gate",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
