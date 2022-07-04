package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeVgames(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(5)

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
		colly.AllowedDomains("store.v-games.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML("li.product.type-product", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price")

		var stock int

		if hasClass(e, "instock") {
			stock = 0
		} else if hasClass(e, "onbackorder") {
			stock = 1
		} else if hasClass(e, "outofstock") {
			stock = 2
		}

		item := models.Price{
			Name:       e.ChildText(".woocommerce-loop-product__title"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".woocommerce-LoopProduct-link.woocommerce-loop-product__link img", "data-src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".woocommerce-LoopProduct-link", "href"),
		}

		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".woocommerce-pagination a.next", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://store.v-games.gr/category/board-games")
	collector.Wait()

	return nil, nil
}
