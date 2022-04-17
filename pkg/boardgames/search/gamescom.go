package search

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/w3m"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func ScrapeGamesCom(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []models.Price{}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	c.OnHTML(".col-tile", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".ty-price")
		name := e.ChildText(".product-title")
		fmt.Println(name, raw_price)
	})

	c.OnHTML("a.ty-pagination__next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		local_link, _ := w3m.Download(link)
		c.Visit(local_link)
	})

	c.OnHTML(".col-tile", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".ty-price")

		rs = append(rs, models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    18,
			StoreThumb: e.ChildAttr(".cm-image", "src"),
			Stock:      childHasClass(e, "button.ty-btn__primary", "ty-btn__add-to-cart"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".product-title", "href"),
		})
	})

	for _, url := range []string{"https://www.gamescom.gr/epitrapezia-el/", "https://www.gamescom.gr/epitrapezia-el/category-124/"} {
		local, err := w3m.Download(url)
		if err != nil {
			return nil, err
		}

		c.Visit(local)
	}

	c.Wait()

	err := updateBatch(db, 18)
	if err != nil {
		return nil, err
	}

	for _, item := range rs {
		err = upsertPrice(db, item)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
