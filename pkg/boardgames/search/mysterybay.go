package search

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"strings"
	"unicode"
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

		if e.ChildText("span[data-hook=product-item-ribbon]") == "PRE-ORDER" {
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

	collector.OnHTML("a.skOBQqy", func(e *colly.HTMLElement) {
		page := strings.Split(e.Attr("data-hook"), "-")

		if len(page) > 1 {
			l, _ := strconv.Atoi(page[1])

			for i := 1; i <= l; i++ {
				link := fmt.Sprintf("%s%d", getPage(e.Request.AbsoluteURL("")), i)
				log.Println("Visiting: " + link)
				collector.Visit(link)
			}
		}
	})

	collector.Visit("https://www.mystery-bay.com/diaxeirisis-poron?page=1")
	collector.Visit("https://www.mystery-bay.com/stratigikis?page=1")
	collector.Visit("https://www.mystery-bay.com/fantasias?page=1")
	collector.Visit("https://www.mystery-bay.com/mystirioy-tromoy?page=1")
	collector.Visit("https://www.mystery-bay.com/paixnidia-me-miniatoyres-dungeon-cr?page=1")
	collector.Visit("https://www.mystery-bay.com/oikogeneiaka?page=1")
	collector.Visit("https://www.mystery-bay.com/tis-pareas?page=1")
	collector.Visit("https://www.mystery-bay.com/paixnidia-me-kartes-zaria?page=1")
	collector.Visit("https://www.mystery-bay.com/lcg?page=1")
	collector.Visit("https://www.mystery-bay.com/war-games?page=1")
	// collector.Visit("https://www.mystery-bay.com/pre-orders")
	collector.Wait()

	return map[string]interface{}{
		"name":     "Mystery Bay",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
}

func getPage(url string) string {
	return strings.TrimRightFunc(url, func(r rune) bool {
		return unicode.IsNumber(r)
	})
}
