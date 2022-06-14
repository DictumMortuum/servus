package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeKaissaEu(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(6)

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
		colly.AllowedDomains("www.kaissa.eu"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML("article.product", func(e *colly.HTMLElement) {
		var stock int
		raw_price := e.ChildText(".price")

		if childHasClass(e, ".add-to-cart input", "stock-update") {
			stock = 2
		} else {
			stock = 0
		}

		item := models.Price{
			Name:       e.ChildText(".caption"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".photo a img", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".photo a", "href")),
		}

		err = insertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".next a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://www.kaissa.eu/products/epitrapezia-kaissa")
	collector.Visit("https://www.kaissa.eu/products/epitrapezia-sta-agglika")
	collector.Wait()

	return nil, nil
}
