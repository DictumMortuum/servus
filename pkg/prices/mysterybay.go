package prices

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
)

func ParseMysteryBay(c *gin.Context) {
	rs, err := createPrices(parseMysteryBay())
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}

func parseMysteryBay() []models.BoardgamePrice {
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
		var data models.BoardgamePrice
		data.Boardgame = e.ChildText("h3")
		data.Store = "Mystery Bay"
		data.OriginalPrice = getPrice(e, "._23IPr")
		data.ReducedPrice = getPrice(e, "._23ArP")

		if data.OriginalPrice-data.ReducedPrice >= PRICE_CUTOFF {
			rs = append(rs, data)
		}
	})

	c.Visit("https://www.mystery-bay.com/prosfores")

	return rs
}
