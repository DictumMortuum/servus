package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

func ScrapeCrystalLotus(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(24)

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
		colly.AllowedDomains("crystallotus.eu"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".grid__item", func(e *colly.HTMLElement) {
		link := e.ChildAttr(".product-card__image-with-placeholder-wrapper img", "data-src")
		if strings.HasPrefix(link, "//") {
			link = "https:" + link
		}

		if strings.Contains(link, "{width}") {
			link = strings.Replace(link, "{width}", "2048", -1)
		}

		raw_price := e.ChildText(".price__sale")
		item := models.Price{
			Name:       e.ChildText("a.grid-view-item__link"),
			StoreId:    store_id,
			StoreThumb: link,
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr("a.grid-view-item__link", "href")),
		}

		err = insertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".pagination a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://crystallotus.eu/collections/board-games")
	collector.Wait()

	return nil, nil
}
