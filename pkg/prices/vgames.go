package prices

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"strings"
)

func ParseVGames(c *gin.Context) {
	rs, err := createPrices(parseVGames())
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}

func parseVGames() []models.BoardgamePrice {
	rs := []models.BoardgamePrice{}

	c := colly.NewCollector(
		colly.AllowedDomains("store.v-games.gr"),
	)

	c.OnHTML("a.next", basicNav(c))

	c.OnHTML("a.woocommerce-LoopProduct-link", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.Index(link, "store.v-games.gr/product") != -1 {
			sale := e.ChildText(".onsale")

			if sale != "" {
				var data models.BoardgamePrice
				data.Boardgame = e.ChildText(".woocommerce-loop-product__title")
				data.Store = "V Games"
				data.OriginalPrice = getPrice(e, ".price > del > .woocommerce-Price-amount")
				data.ReducedPrice = getPrice(e, ".price > ins > .woocommerce-Price-amount")

				if data.OriginalPrice-data.ReducedPrice >= PRICE_CUTOFF {
					rs = append(rs, data)
				}
			}
		}
	})

	c.Visit("https://store.v-games.gr")

	return rs
}
