package scraper

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
)

type VgamesScraper struct {
	Store string
}

func (obj VgamesScraper) GetStore() models.Store {
	return models.Store{0, obj.Store}
}

func (obj VgamesScraper) Scrape() ([]models.ScraperData, error) {
	rs := []models.ScraperData{}

	database, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	id, err := db.Exists(database, obj.GetStore())
	if err != nil {
		return nil, err
	}

	c := colly.NewCollector(
		colly.AllowedDomains("store.v-games.gr"),
	)

	c.OnHTML("a.next", basicNav(c))

	c.OnHTML(".product.type-product ", func(e *colly.HTMLElement) {
		var data models.ScraperData
		data.StoreId = *id
		data.Link = e.ChildAttr(".woocommerce-LoopProduct-link", "href")
		data.Title = e.ChildText(".woocommerce-loop-product__title")
		data.SKU = e.ChildAttr(".button", "data-product_sku")
		rs = append(rs, data)
	})

	c.Visit("https://store.v-games.gr/category/board-games")

	inserted_rs := []models.ScraperData{}

	for _, item := range rs {
		id, err := db.InsertIfNotExists(database, item)
		if err != nil {
			return nil, err
		}

		if id != nil {
			inserted_rs = append(inserted_rs, item)
		}
	}

	return inserted_rs, nil
}

func (obj VgamesScraper) ScrapePrices() ([]models.ScraperPrice, error) {
	rs := []models.ScraperPrice{}

	database, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	links, err := GetDataWithLeastNoPrices(database, obj)
	if err != nil {
		return nil, err
	}

	for _, item := range links {
		c := colly.NewCollector(
			colly.AllowedDomains("store.v-games.gr"),
		)

		c.OnHTML(".summary", func(e *colly.HTMLElement) {
			var data models.ScraperPrice
			data.DataId = item.Id

			price := e.ChildText(".price > ins > .woocommerce-Price-amount")

			if price == "" {
				price = e.ChildText(".price > .woocommerce-Price-amount")
			}

			data.Price = getPriceString(price)
			data.InStock = e.ChildText(".out-of-stock") != ""
			rs = append(rs, data)
		})

		c.Visit(item.Link)
	}

	for _, item := range rs {
		_, err := db.InsertIfNotExists(database, item)
		if err != nil {
			return nil, err
		}
	}

	return rs, nil
}
