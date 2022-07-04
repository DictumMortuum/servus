package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeEfantasy(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(8)

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
		colly.AllowedDomains("www.efantasy.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML("div.product.product-box", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".product-price a strong")

		var stock int

		if e.Attr("data-label") == "" {
			stock = 2
		} else if e.Attr("data-label") == "preorder" {
			stock = 1
		} else {
			stock = 0
		}

		item := models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".product-image a img", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".product-title a", "href")),
		}

		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".pagination a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://www.efantasy.gr/en/products/%CE%B5%CF%80%CE%B9%CF%84%CF%81%CE%B1%CF%80%CE%AD%CE%B6%CE%B9%CE%B1-%CF%80%CE%B1%CE%B9%CF%87%CE%BD%CE%AF%CE%B4%CE%B9%CE%B1/sc-all")
	collector.Wait()

	return nil, nil
}
