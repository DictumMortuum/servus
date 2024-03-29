package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeRollnplay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(26)
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
		colly.AllowedDomains("rollnplay.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".product.type-product", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price ins .amount")

		if raw_price == "" {
			raw_price = e.ChildText(".price .amount")
		}

		var stock int

		if hasClass(e, "instock") {
			stock = 0
		} else if hasClass(e, "onbackorder") {
			stock = 1
		} else if hasClass(e, "outofstock") {
			stock = 2
		}

		item := models.Price{
			Name:       e.ChildText(".heading-title"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".has-back-image img", "data-src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".heading-title a", "href"),
		}

		detected++
		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".woocommerce-pagination a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://rollnplay.gr/?product_cat=all-categories")
	collector.Wait()

	return map[string]interface{}{
		"name":     "Roll n Play",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
