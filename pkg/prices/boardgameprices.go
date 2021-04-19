package prices

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"regexp"
	"strconv"
)

func ParseBoardgameprices(c *gin.Context) {
	rs, err := createPrices(parseBoardgameprices())
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}

func parseBoardgameprices() []models.BoardgamePrice {
	rs := []models.BoardgamePrice{}
	re_from := regexp.MustCompile("from €([0-9]+.[0-9]+)")
	re_to := regexp.MustCompile("to €([0-9]+.[0-9]+)")

	c := colly.NewCollector(
		colly.AllowedDomains("boardgameprices.co.uk"),
	)

	c.OnHTML("#searchresultlist .searchinfocontainer.multicell", func(e *colly.HTMLElement) {
		var data models.BoardgamePrice
		data.Boardgame = e.ChildText(".searchcell.itemname span a")
		data.Store = e.ChildText(".searchcell .storename")
		raw := e.ChildText(".searchcell")

		from := re_from.FindStringSubmatch(raw)
		to := re_to.FindStringSubmatch(raw)

		if len(from) > 0 {
			data.OriginalPrice, _ = strconv.ParseFloat(from[1], 64)
		}

		if len(to) > 0 {
			data.ReducedPrice, _ = strconv.ParseFloat(to[1], 64)
		}

		rs = append(rs, data)
	})

	c.Visit("https://boardgameprices.co.uk/item/pricedrops?order=date&country=GR&minimum=10")

	return rs
}
