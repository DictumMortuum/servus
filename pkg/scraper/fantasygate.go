package scraper

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
)

type FantasyGateScraper struct {
	Store string
}

func (obj FantasyGateScraper) GetStore() models.Store {
	return models.Store{0, obj.Store}
}

func (obj FantasyGateScraper) Scrape() ([]models.ScraperData, error) {
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
		colly.AllowedDomains("fantasygate.gr"),
	)

	c.OnHTML("a.pagenav", basicNav(c))

	c.OnHTML(".block_product", func(e *colly.HTMLElement) {
		var data models.ScraperData
		data.StoreId = *id
		data.Link = e.Request.AbsoluteURL(e.ChildAttr(".name a", "href"))
		data.Title = e.ChildText(".name")
		data.SKU = ""
		rs = append(rs, data)
	})

	c.Visit("https://fantasygate.gr/family-games")
	c.Visit("https://fantasygate.gr/fantasygames")
	c.Visit("https://fantasygate.gr/strategygames")
	c.Visit("https://fantasygate.gr/cardgames")

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

func (obj FantasyGateScraper) ScrapePrices() ([]models.ScraperPrice, error) {
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
			colly.AllowedDomains("fantasygate.gr"),
		)

		c.OnHTML("#comjshop > form", func(e *colly.HTMLElement) {
			var data models.ScraperPrice
			data.DataId = item.Id
			price := e.ChildText("#block_price")
			data.Price = getPriceString(price)
			data.InStock = e.ChildText("#not_available") == ""
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
