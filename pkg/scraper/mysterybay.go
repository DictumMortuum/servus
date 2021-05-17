package scraper

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
)

type MysteryBayScraper struct {
	Store string
}

func (obj MysteryBayScraper) GetStore() models.Store {
	return models.Store{0, obj.Store}
}

func (obj MysteryBayScraper) Scrape() ([]models.ScraperData, error) {
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
		colly.AllowedDomains("www.mystery-bay.com"),
	)

	c.OnHTML("span", func(e *colly.HTMLElement) {
		hook := e.Attr("data-hook")
		pagination := strings.Split(hook, " ")

		if len(pagination) == 2 {
			if pagination[1] == "current-page" {
				page := strings.Split(pagination[0], "-")
				next, _ := strconv.ParseInt(page[1], 10, 64)
				c.Visit(fmt.Sprintf("https://www.mystery-bay.com/ola-ta-proionta?page=%d", next+1))
			}
		}
	})

	c.OnHTML("div._2zTHN._2AHc6", func(e *colly.HTMLElement) {
		var data models.ScraperData
		data.StoreId = *id
		data.Link = e.ChildAttr("._34sIs", "href")
		data.Title = e.ChildText("._3RqKm h3")
		data.SKU = ""
		rs = append(rs, data)
	})

	c.Visit("https://www.mystery-bay.com/ola-ta-proionta")

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

func (obj MysteryBayScraper) ScrapePrices() ([]models.ScraperPrice, error) {
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
			colly.AllowedDomains("www.mystery-bay.com"),
		)

		c.OnHTML("._12vNY", func(e *colly.HTMLElement) {
			var data models.ScraperPrice
			data.DataId = item.Id
			price := e.ChildText("._26qxh > span:not(._19Hjy)")
			data.Price = getPriceString(price)
			data.InStock = e.ChildAttr("._3j0qu > button", "aria-disabled") == "true"
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
