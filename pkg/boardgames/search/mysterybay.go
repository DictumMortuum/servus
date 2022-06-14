package search

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

func ScrapeMysteryBay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(3)

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
		colly.AllowedDomains("www.mystery-bay.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML("._3DNsL", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "") {
			return
		}

		raw_price := e.ChildText("._2-l9W")
		raw_url := e.ChildAttr("._1FMIK", "style")
		urls := getURL(raw_url)

		url := ""
		if len(urls) > 0 {
			url = urls[0]
		}

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

		err = insertQueueItem(ch, q, item)
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

	return nil, nil
}
