package scraper

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	price = regexp.MustCompile("([0-9]+.[0-9]+)")
)

func getPriceString(raw string) float64 {
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

func Scrape(obj models.Scrapable) func(*gin.Context) {
	return func(c *gin.Context) {
		rs, err := obj.Scrape()
		if err != nil {
			util.Error(c, err)
			return
		}

		util.Success(c, &rs)
	}
}

func ScrapePrices(obj models.Scrapable) func(*gin.Context) {
	return func(c *gin.Context) {
		rs, err := obj.ScrapePrices()
		if err != nil {
			util.Error(c, err)
			return
		}

		util.Success(c, &rs)
	}
}

func GetDataWithLeastNoPrices(database *sqlx.DB, obj models.Scrapable) ([]models.ScraperData, error) {
	rs := []models.ScraperData{}

	id, err := db.Exists(database, obj.GetStore())
	if err != nil {
		return nil, err
	}

	udb := database.Unsafe()
	err = udb.Select(&rs, `
	select
		d.*,
		count(*) cnt
	from
		tboardgamescraperdata d
	left join
		tboardgamescraperprices p
	on
		d.id = p.data_id
	where
		d.store_id = ?
	group by
		d.id
	order by
		cnt
	`, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func CreateDataMapping(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)
	var game models.ScraperData

	log.Println(data)

	if val, ok := data["id"]; ok {
		game.Id = int64(val.(float64))
	} else {
		util.Error(c, errors.New("asdf"))
		return
	}

	if val, ok := data["boardgame_id"]; ok {
		game.BoardgameId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	} else {
		util.Error(c, errors.New("asdf"))
		return
	}

	database, err := db.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	sql := `update tboardgamescraperdata set boardgame_id = :boardgame_id where id = :id`
	_, err = database.NamedExec(sql, &game)
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, game)
}
