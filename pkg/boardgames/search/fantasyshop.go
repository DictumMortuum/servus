package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/w3m"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func ScrapeFantasyShop(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(28)

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

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	collector := colly.NewCollector()
	collector.WithTransport(t)

	collector.OnHTML(".ty-column3", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".ty-price-num")

		item := models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".ty-pict.cm-image", "src"),
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".ty-grid-list__image a", "href"),
		}

		err = insertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML("a.ty-pagination__next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println("Visiting: " + link)
		local_link, _ := w3m.BypassCloudflare(link)
		collector.Visit(local_link)
	})

	local, err := w3m.BypassCloudflare("https://www.fantasy-shop.gr/epitrapezia-paihnidia-pazl/?features_hash=18-Y")
	if err != nil {
		return nil, err
	}

	collector.Visit(local)
	collector.Wait()

	return nil, nil
}
