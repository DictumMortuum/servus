package search

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeAvalon(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(25)

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
		colly.AllowedDomains("avalongames.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".product-layout", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price-normal")

		if raw_price == "" {
			raw_price = e.ChildText(".price-new")
		}

		var stock int

		if !hasClass(e, ".out-of-stock") {
			stock = 0
		} else {
			stock = 2
		}

		item := models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".product-img div img", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
		}

		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".pagination-results", func(e *colly.HTMLElement) {
		pageCount := getPages(e.Text)
		for i := 2; i <= pageCount; i++ {
			link := fmt.Sprintf("https://avalongames.gr/index.php?route=product/category&path=59&limit=100&page=%d", i)
			log.Println("Visiting: ", link)
			collector.Visit(link)
		}
	})

	collector.Visit("https://avalongames.gr/index.php?route=product/category&path=59&limit=100")
	collector.Wait()

	return nil, nil
}
