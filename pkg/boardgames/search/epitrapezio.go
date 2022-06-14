package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeEpitrapezio(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(15)

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
		colly.AllowedDomains("epitrapez.io"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML("li.product.type-product", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price ins .amount")

		if raw_price == "" {
			raw_price = e.ChildText(".price .amount")
		}

		var stock int

		if e.ChildText("a.add_to_cart_button") != "" {
			stock = 0
		} else {
			stock = 2
		}

		item := models.Price{
			Name:       e.ChildText(".woocommerce-loop-product__title"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".epz-product-thumbnail img", "data-src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".woocommerce-LoopProduct-link", "href"),
		}

		err = insertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".woocommerce-pagination a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://epitrapez.io/product-category/epitrapezia/?Stock=allstock")
	collector.Wait()

	return nil, nil
}
