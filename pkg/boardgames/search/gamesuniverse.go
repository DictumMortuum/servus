package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeGamesUniverse(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(20)

	log.Printf("Scraper %d started\n", store_id)

	conn, ch, q, err := setupQueue("prices")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer ch.Close()

	err = updateBatch(db, store_id)
	if err != nil {
		return nil, err
	}

	collector := colly.NewCollector(
		colly.AllowedDomains("gamesuniverse.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML("article.product-miniature", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".product-price")

		item := models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".thumbnail img", "data-src"),
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".product-thumbnail", "href"),
		}

		err = insertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://gamesuniverse.gr/el/10-epitrapezia")
	collector.Wait()

	return nil, nil
}
