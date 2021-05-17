package scraper

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
)

type KaissaScraper struct {
	Store string
}

func (obj KaissaScraper) GetStore() models.Store {
	return models.Store{0, obj.Store}
}

func (obj KaissaScraper) Scrape() ([]models.ScraperData, error) {
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
		colly.AllowedDomains("www.kaissa.eu"),
	)

	c.OnHTML(".next a", basicNav(c))

	c.OnHTML("article.product", func(e *colly.HTMLElement) {
		var data models.ScraperData
		data.StoreId = *id
		data.Link = e.Request.AbsoluteURL(e.ChildAttr(".photo a", "href"))
		data.Title = e.ChildText(".caption")
		data.SKU = ""
		rs = append(rs, data)
	})

	c.Visit("https://www.kaissa.eu/products/epitrapezia-kaissa")
	c.Visit("https://www.kaissa.eu/products/epitrapezia-sta-agglika")

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

func (obj KaissaScraper) ScrapePrices() ([]models.ScraperPrice, error) {
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
			colly.AllowedDomains("www.kaissa.eu"),
		)

		c.OnHTML(".product-details", func(e *colly.HTMLElement) {
			var data models.ScraperPrice
			data.DataId = item.Id
			price := e.ChildText(".price .final")
			data.Price = getPriceString(price)
			data.InStock = e.ChildText("#product-subscribe-form") == ""
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
