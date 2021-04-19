package prices

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"strings"
)

func ParseFantasyGate(c *gin.Context) {
	rs, err := createPrices(parseFantasyGate())
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}

func parseFantasyGate() []models.BoardgamePrice {
	rs := []models.BoardgamePrice{}

	c := colly.NewCollector(
		colly.AllowedDomains("fantasygate.gr"),
	)

	detailCollector := c.Clone()

	c.OnHTML("a.pagenav", basicNav(c))

	c.OnHTML(".sblock4", func(e *colly.HTMLElement) {
		link := e.ChildAttr(".product > .name > a", "href")
		link = e.Request.AbsoluteURL(link)

		from := getPrice(e, ".old_price > span")
		to := getPrice(e, ".jshop_price > span")

		if from > 0.0 {
			if from-to >= PRICE_CUTOFF {
				detailCollector.Visit(link)
			}
		}
	})

	detailCollector.OnHTML(`#comjshop`, func(e *colly.HTMLElement) {
		var data models.BoardgamePrice
		data.Boardgame = strings.TrimSuffix(e.ChildText("h1"), " "+e.ChildText(".jshop_code_prod"))
		data.Store = "Fantasy Gate"
		data.OriginalPrice = getPrice(e, "#old_price")
		data.ReducedPrice = getPrice(e, "#block_price")
		rs = append(rs, data)
	})

	c.Visit("https://fantasygate.gr/family-games")
	c.Visit("https://fantasygate.gr/fantasygames")
	c.Visit("https://fantasygate.gr/strategygames")
	c.Visit("https://fantasygate.gr/cardgames")

	return rs
}
