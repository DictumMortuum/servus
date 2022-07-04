package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeMeeplePlanet(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(7)

	log.Printf("Scraper %d started\n", store_id)

	conn, ch, q, err := rabbitmq.SetupQueue("prices")
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
		colly.AllowedDomains("meeple-planet.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML("div.product-small.product-type-simple", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".amount")

		var stock int

		if hasClass(e, "out-of-stock") {
			stock = 2
		} else {
			if e.ChildText(".badge-inner") != "" {
				stock = 1
			} else {
				stock = 0
			}
		}

		item := models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		}

		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML("a.next.page-number", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://meeple-planet.com/category/epitrapezia-paixnidia")
	collector.Visit("https://meeple-planet.com/category/pre-orders")
	collector.Wait()

	return nil, nil
}
