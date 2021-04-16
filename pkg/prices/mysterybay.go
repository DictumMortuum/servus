package prices

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
)

func ParseMysteryBay() []models.BoardgamePrice {
	rs := []models.BoardgamePrice{}

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
				c.Visit(fmt.Sprintf("https://www.mystery-bay.com/prosfores?page=%d", next+1))
			}
		}
	})

	c.OnHTML("div._2zTHN._2AHc6", func(e *colly.HTMLElement) {
		from := getPrice(e, "._23IPr")
		to := getPrice(e, "._23ArP")

		if from > 0.0 && from-to >= PRICE_CUTOFF {
			var data models.BoardgamePrice
			data.Boardgame = e.ChildText("h3")
			data.Store = "Mystery Bay"
			data.OriginalPrice = from
			data.ReducedPrice = to
			data.PriceDiff = data.OriginalPrice - data.ReducedPrice
			//e.ChildText(".sku")
			rs = append(rs, data)
		}
	})

	c.Visit("https://www.mystery-bay.com/prosfores")

	return rs
}
