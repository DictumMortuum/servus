package prices

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
)

func ParseGameRules(c *gin.Context) {
	rs, err := createPrices(parseGameRules())
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}

func parseGameRules() []models.BoardgamePrice {
	rs := []models.BoardgamePrice{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.thegamerules.com"),
	)

	detailCollector := c.Clone()

	c.OnHTML("a.next", basicNav(c))

	c.OnHTML("a.product-img.has-second-image", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.Index(link, "thegamerules.com/offers") != -1 {
			detailCollector.Visit(link)
		}
	})

	detailCollector.OnHTML(`#content`, func(e *colly.HTMLElement) {
		var data models.BoardgamePrice
		data.Boardgame = e.ChildText("div.title.page-title")
		data.Store = "The Game Rules"

		from := price.FindStringSubmatch(e.ChildText(".product-price-old"))
		to := price.FindStringSubmatch(e.ChildText(".product-price-new"))

		if len(from) > 0 {
			data.OriginalPrice, _ = strconv.ParseFloat(from[1], 64)
		}

		if len(to) > 0 {
			data.ReducedPrice, _ = strconv.ParseFloat(to[1], 64)
		}

		if data.OriginalPrice-data.ReducedPrice >= PRICE_CUTOFF {
			rs = append(rs, data)
		}
	})

	c.Visit("https://www.thegamerules.com/offers?fa132=Board%20Game%20Expansions,Board%20Games")

	return rs
}
