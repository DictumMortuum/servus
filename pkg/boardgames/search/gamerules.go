package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapeGameRules(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(4)
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
		colly.AllowedDomains("www.thegamerules.com"),
		colly.CacheDir("/tmp/servus"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".product-layout", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price-new")

		if raw_price == "" {
			raw_price = e.ChildText(".price-normal")
		}

		var stock int

		switch e.ChildText(".c--stock-label") {
		case "Εκτός αποθέματος":
			stock = 2
		case "Άμεσα Διαθέσιμο":
			stock = 0
		default:
			stock = 1
		}

		item := models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".product-img div img", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		}

		detected++
		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://www.thegamerules.com/epitrapezia-paixnidia?fa132=Board%20Game%20Expansions")
	collector.Visit("https://www.thegamerules.com/epitrapezia-paixnidia?fa132=Board%20Games")
	collector.Wait()

	return map[string]interface{}{
		"name":     "Game Rules",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
