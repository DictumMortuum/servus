package prices

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"strings"
)

func ParseKaissa(c *gin.Context) {
	rs, err := createPrices(parseKaissa())
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}

func parseKaissa() []models.BoardgamePrice {
	rs := []models.BoardgamePrice{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.kaissa.eu"),
	)

	c.OnHTML(".next a", basicNav(c))

	c.OnHTML(".product", func(e *colly.HTMLElement) {
		var data models.BoardgamePrice
		data.Boardgame = e.ChildText(".product .caption a")
		data.Store = "Kaissa Amarousiou"

		from := e.ChildText(".original-price")
		to := e.ChildText(".price")
		to = strings.TrimPrefix(to, from)

		data.OriginalPrice = getPriceString(from)
		data.ReducedPrice = getPriceString(to)

		if data.OriginalPrice-data.ReducedPrice >= PRICE_CUTOFF {
			rs = append(rs, data)
		}
	})

	c.Visit("https://www.kaissa.eu/products/offers")

	return rs
}
