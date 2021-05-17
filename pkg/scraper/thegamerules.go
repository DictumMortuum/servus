package scraper

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
)

type GameRulesScraper struct {
	Store string
}

func (obj GameRulesScraper) GetStore() models.Store {
	return models.Store{0, obj.Store}
}

func (obj GameRulesScraper) Scrape() ([]models.ScraperData, error) {
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
		colly.AllowedDomains("www.thegamerules.com"),
	)

	c.OnHTML("a.next", basicNav(c))

	c.OnHTML(".product-layout", func(e *colly.HTMLElement) {
		var data models.ScraperData
		data.StoreId = *id
		data.Link = e.ChildAttr(".product-img", "href")
		data.Title = e.ChildText(".name")
		data.SKU = ""
		rs = append(rs, data)
	})

	c.Visit("https://www.thegamerules.com/epitrapezia-paixnidia?fa132=Board%20Games")

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

func (obj GameRulesScraper) ScrapePrices() ([]models.ScraperPrice, error) {
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
			colly.AllowedDomains("www.thegamerules.com"),
		)

		c.OnHTML("#product", func(e *colly.HTMLElement) {
			var data models.ScraperPrice
			data.DataId = item.Id

			price := e.ChildText(".product-price-new")

			if price == "" {
				price = e.ChildText(".product-price")
			}

			data.Price = getPriceString(price)
			data.InStock = e.ChildText(".product-stock") == "Άμεσα Διαθέσιμο"
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
