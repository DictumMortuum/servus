package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

func ScrapePoliteia(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(12)

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

	collector := colly.NewCollector(
		colly.AllowedDomains("www.politeianet.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".browse-page-block", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".productPrice")
		if raw_price == "" {
			return
		}

		item := models.Price{
			Name:       e.ChildText(".browse-product-title"),
			StoreId:    store_id,
			StoreThumb: e.ChildAttr(".browseProductImage", "src"),
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".browse-product-title", "href"),
		}

		err = rabbitmq.InsertQueueItem(ch, q, item)
		if err != nil {
			log.Println(err)
		}
	})

	collector.OnHTML("a.pagenav", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		collector.Visit(link)
	})

	collector.Visit("https://www.politeianet.gr/index.php?option=com_virtuemart&category_id=948&page=shop.browse&subCatFilter=-1&langFilter=-1&pubdateFilter=-1&availabilityFilter=-1&discountFilter=-1&priceFilter=-1&pageFilter=-1&Itemid=721&limit=20&limitstart=0")
	collector.Wait()

	return nil, nil
}
