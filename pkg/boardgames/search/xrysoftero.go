package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeXrysoFtero(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(21)

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
		colly.AllowedDomains("xrysoftero.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".thumbnail-container", func(e *colly.HTMLElement) {
		url := e.ChildAttr(".cover-image", "src")
		if url == "" {
			return
		}

		raw_price := e.ChildText(".price")
		item := models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    store_id,
			StoreThumb: url,
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr("a.relative", "href"),
		}

		err = insertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://xrysoftero.gr/362-epitrapezia-paixnidia?resultsPerPage=48&q=%CE%9C%CE%AC%CF%81%CE%BA%CE%B1%5C-%CE%95%CE%BA%CE%B4%CF%8C%CF%84%CE%B7%CF%82-%CE%9A%CE%AC%CE%B9%CF%83%CF%83%CE%B1")
	collector.Wait()

	return nil, nil
}
