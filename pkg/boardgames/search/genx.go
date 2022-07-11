package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeGenx(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(27)
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
		colly.AllowedDomains("www.genx.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".white_bg", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".txtSale")

		if raw_price == "" {
			raw_price = e.ChildText(".txtPrice")
		}

		raw_stock := e.ChildText(".txtOutOfStock")

		var stock int

		if raw_stock == "" {
			stock = 0
		} else {
			stock = 2
		}

		item := models.Price{
			Name:       e.ChildText(".txtTitle"),
			StoreId:    store_id,
			StoreThumb: e.Request.AbsoluteURL(e.ChildAttr(".hover01 a img", "src")),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".hover01 a", "href")),
		}

		detected++
		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".prevnext", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://www.genx.gr/index.php?page=0&act=viewCat&catId=60&prdsPage=45")
	collector.Wait()

	return map[string]interface{}{
		"name":     "Genx",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
