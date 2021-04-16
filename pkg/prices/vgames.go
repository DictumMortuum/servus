package prices

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	// "regexp"
	// "fmt"
	"strconv"
	"strings"
)

func getPrice(e *colly.HTMLElement, selector string) float64 {
	raw := e.ChildText(selector)
	raw = strings.ReplaceAll(raw, ",", ".")
	match := price.FindStringSubmatch(raw)

	if len(match) > 0 {
		price, _ := strconv.ParseFloat(match[1], 64)
		return price
	} else {
		return 0.0
	}
}

func basicNav(c *colly.Collector) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	}
}

func ParseVGames() []models.BoardgamePrice {
	rs := []models.BoardgamePrice{}

	c := colly.NewCollector(
		colly.AllowedDomains("store.v-games.gr"),
	)

	detailCollector := c.Clone()

	c.OnHTML("a.next", basicNav(c))

	c.OnHTML("a.woocommerce-LoopProduct-link", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.Index(link, "store.v-games.gr/product") != -1 {
			sale := e.ChildText(".onsale")

			if sale != "" {
				from := getPrice(e, ".price > del > .woocommerce-Price-amount")
				to := getPrice(e, ".price > ins > .woocommerce-Price-amount")

				if from-to >= PRICE_CUTOFF {
					// fmt.Printf("Offer found: %s\n", link)
					detailCollector.Visit(link)
				}
			}
		}
	})

	detailCollector.OnHTML(`#primary`, func(e *colly.HTMLElement) {
		var data models.BoardgamePrice
		data.Boardgame = e.ChildText("h1.product_title")
		data.Store = "V Games"
		data.OriginalPrice = getPrice(e, ".price > del > .woocommerce-Price-amount")
		data.ReducedPrice = getPrice(e, ".price > ins > .woocommerce-Price-amount")
		data.PriceDiff = data.OriginalPrice - data.ReducedPrice
		//e.ChildText(".sku")
		rs = append(rs, data)
	})

	c.Visit("https://store.v-games.gr")

	return rs
}
