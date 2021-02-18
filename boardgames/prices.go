package boardgames

import (
	"errors"
	"fmt"
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetPrices(c *gin.Context) {
	order := c.DefaultQuery("order", "date")
	country := c.DefaultQuery("country", "GR")
	min := c.DefaultQuery("minimum", "10")
	retval := []PriceRow{}

	re_from := regexp.MustCompile("from €([0-9]+.[0-9]+)")
	re_to := regexp.MustCompile("to €([0-9]+.[0-9]+)")
	re_date := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}")

	res, err := http.Get("https://boardgameprices.co.uk/item/pricedrops?order=" + order + "&country=" + country + "&minimum=" + min)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		util.Error(c, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)))
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		util.Error(c, err)
		return
	}

	database, err := db.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	doc.Find("#searchresultlist .searchinfocontainer.multicell").Each(func(i int, s *goquery.Selection) {
		var data PriceRow
		data.Boardgame = s.Find(".searchcell .name").Text()
		data.Store = strings.TrimSpace(s.Find(".searchcell .storename").Text())
		raw := s.Text()

		from := re_from.FindStringSubmatch(raw)
		to := re_to.FindStringSubmatch(raw)
		date := re_date.FindStringSubmatch(raw)

		if len(from) > 0 {
			data.OriginalPrice, _ = strconv.ParseFloat(from[1], 64)
		}

		if len(to) > 0 {
			data.ReducedPrice, _ = strconv.ParseFloat(to[1], 64)
		}

		if len(date) > 0 {
			data.Date, _ = time.Parse("2006-01-02", date[0])
		}

		if len(from) > 0 && len(to) > 0 {
			data.PriceDiff = data.OriginalPrice - data.ReducedPrice
		}

		id, err := PriceExists(database, data)
		if err != nil {
			util.Error(c, err)
			return
		}

		data.Id = id

		if id > 0 {
			err = UpdatePrice(database, data)
			if err != nil {
				util.Error(c, err)
				return
			}
		} else {
			err = CreatePrice(database, data)
			if err != nil {
				util.Error(c, err)
				return
			}

			msg := fmt.Sprintf("%s offers %s at %.2f from %.2f\n", data.Store, data.Boardgame, data.ReducedPrice, data.OriginalPrice)
			err = util.TelegramMessage(msg)
			if err != nil {
				util.Error(c, err)
				return
			}
		}

		retval = append(retval, data)
	})

	util.Success(c, retval)
}
