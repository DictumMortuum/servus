package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeGameExplorers(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(22)
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
		colly.AllowedDomains("www.gameexplorers.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".single-product-item", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".regular-price")
		item := models.Price{
			Name:       e.ChildText("h2:nth-child(1)"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr("a:nth-child(1) > img:nth-child(1)", "src"),
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr("a:nth-child(1)", "href"),
		}

		detected++
		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".product-pagination > a", func(e *colly.HTMLElement) {
		if e.Attr("title") == "επόμενη σελίδα" {
			link := e.Attr("href")
			log.Println("Visiting: " + link)
			collector.Visit(link)
		}
	})

	collector.Visit("https://www.gameexplorers.gr/kartes-epitrapezia/epitrapezia-paixnidia/items-grid-date-desc-1-60.html")
	collector.Wait()

	return map[string]interface{}{
		"name":     "Game Explorers",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
