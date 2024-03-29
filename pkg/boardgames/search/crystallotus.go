package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/DictumMortuum/servus/pkg/w3m"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strings"
)

func ScrapeCrystalLotus(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(24)
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

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	collector := colly.NewCollector()
	collector.WithTransport(t)

	collector.OnHTML(".grid__item", func(e *colly.HTMLElement) {
		link := e.ChildAttr(".motion-reduce", "src")
		if strings.HasPrefix(link, "//") {
			link = "https:" + link
		}

		raw_price := e.ChildText(".price__sale")
		item := models.Price{
			Name:       e.ChildText(".card-information__text"),
			StoreId:    store_id,
			StoreThumb: link,
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        "https://crystallotus.eu" + e.ChildAttr("a.card-information__text", "href"),
		}

		detected++
		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML(".pagination__list li:last-child a", func(e *colly.HTMLElement) {
		link := "https://crystallotus.eu" + e.Attr("href")
		log.Println("Visiting: " + link)
		local_link, _ := w3m.BypassCloudflare(link)
		collector.Visit(local_link)
	})

	local, err := w3m.BypassCloudflare("https://crystallotus.eu/collections/board-games")
	if err != nil {
		return nil, err
	}

	collector.Visit(local)
	collector.Wait()

	return map[string]interface{}{
		"name":     "Crystal Lotus",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
