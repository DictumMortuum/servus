package search

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

func ScrapeMysteryBay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(3)
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
		colly.AllowedDomains("www.mystery-bay.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML("li[data-hook=product-list-grid-item]", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "") {
			return
		}

		raw_price := e.ChildText("span[data-hook=product-item-price-to-pay]")
		raw_url := e.ChildAttr("[data-hook=product-item-images]", "style")
		urls := getURL(raw_url)

		url := ""
		if len(urls) > 0 {
			candidate := urls[0]
			filtered := strings.Split(candidate, "/v1/fill")
			url = filtered[0]
		}
		// style="background-image:url(https://static.wixstatic.com/media/9dcd7c_df5e66ff7168447ab10021bfa739a4cc~mv2.png/v1/fill/w_100,h_100,al_c,usm_0.66_1.00_0.01/9dcd7c_df5e66ff7168447ab10021bfa739a4cc~mv2.png);background-size:contain" data-hook="">
		var stock int

		if e.ChildText("span[data-hook=product-item-ribbon-new]") == "PRE-ORDER" {
			stock = 1
		} else {
			if e.ChildAttr("button[data-hook=product-item-add-to-cart-button]", "aria-disabled") == "true" {
				stock = 2
			} else {
				stock = 0
			}
		}

		item := models.Price{
			Name:       e.ChildText("h3"),
			StoreId:    store_id,
			StoreThumb: url,
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr("a[data-hook=product-item-container]", "href"),
		}

		detected++
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

	collector.Visit("https://www.mystery-bay.com/epitrapezia-paixnidia?page=36")
	collector.Wait()

	return map[string]interface{}{
		"name":     "Mystery Bay",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}
