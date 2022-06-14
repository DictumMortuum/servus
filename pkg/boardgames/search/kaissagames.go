package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeKaissaGames(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(9)

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

	collector.Visit("https://kaissagames.com/b2c_gr/xenoglossa-epitrapezia.html")
	collector.Wait()

	return nil, nil
}
