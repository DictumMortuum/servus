package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/DictumMortuum/servus/pkg/w3m"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func ScrapeMeepleOnBoard(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(10)

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

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	collector := colly.NewCollector()
	collector.WithTransport(t)

	collector.OnHTML("div.product-small.purchasable", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".amount")

		var stock int

		if hasClass(e, "instock") {
			stock = 0
		} else if hasClass(e, "onbackorder") {
			stock = 1
		} else if hasClass(e, "out-of-stock") {
			stock = 2
		}

		item := models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
		}

		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println("Visiting: " + link)
		local_link, _ := w3m.BypassCloudflare(link)
		collector.Visit(local_link)
	})

	local, err := w3m.BypassCloudflare("https://meepleonboard.gr/product-category/board-games")
	if err != nil {
		return nil, err
	}

	collector.Visit(local)
	collector.Wait()

	return nil, nil
}
