package prices

import (
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/telegram"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"regexp"
	"strconv"
	"strings"
)

var (
	PRICE_CUTOFF = 10.0
	price        = regexp.MustCompile("([0-9]+.[0-9]+)")
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

func GetUpdates(c *gin.Context) {
	database, err := db.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	raw, err := telegram.GetUpdates(database)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Success(c, raw)
}

func SendMessages(c *gin.Context) {
	rs, err := sendMessages()
	if err != nil {
		util.Error(c, err)
	}

	util.Success(c, &rs)
}

func sendMessages() ([]models.Message, error) {
	var data []models.Message
	rs := []models.Message{}

	database, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	err = database.Select(&data, "select * from tmsg where date_send is NULL")
	if err != nil {
		return nil, err
	}

	for _, message := range data {
		err := telegram.TelegramMessage(database, message.Msg)
		if err != nil {
			return nil, err
		}

		_, err = database.NamedExec("update tmsg set date_send = NOW() where id = :id", message)
		if err != nil {
			return nil, err
		}

		rs = append(rs, message)
	}

	return rs, nil
}

func CreateMessages(c *gin.Context) {
	rs, err := createMessages()
	if err != nil {
		util.Error(c, err)
	}

	util.Success(c, &rs)
}

func createMessages() ([]models.BoardgamePrice, error) {
	var data []models.BoardgamePrice
	rs := []models.BoardgamePrice{}

	database, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	err = database.Select(&data, `
	select
		p.id,
		p.cr_date,
		p.boardgame_id,
		g.name as boardgame,
		p.store_id,
		s.name as store,
		p.original_price,
		p.reduced_price
	from
		tboardgameprices p,
		tboardgames g,
		tboardgamestores s
	where
		p.boardgame_id = g.id and
		p.store_id = s.id and
		g.unmapped = 0
	`)
	if err != nil {
		return nil, err
	}

	for _, price := range data {
		inserted, err := db.InsertMsg(database, price)
		if err != nil {
			return nil, err
		}

		if inserted {
			rs = append(rs, price)
		}
	}

	return rs, nil
}

func createPrices(data []models.BoardgamePrice) ([]models.BoardgamePrice, error) {
	rs := []models.BoardgamePrice{}

	database, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	for _, price := range data {
		game := models.Boardgame{
			Name: price.Boardgame,
		}

		atlas, err := boardgames.AtlasSearch(game)
		if err != nil {
			return nil, err
		}
		if len(atlas) == 0 {
			// do not create the game if there are no atlas entries
			continue
		}

		id, err := db.InsertIfNotExists(database, game)
		if err != nil {
			return nil, err
		}

		price.BoardgameId = *id

		store := models.Store{
			Name: price.Store,
		}

		id, err = db.Exists(database, store)
		if err != nil {
			return nil, err
		}

		price.StoreId = *id

		_, err = db.InsertIfNotExists(database, price)
		if err != nil {
			return nil, err
		}

		rs = append(rs, price)
	}

	return rs, nil
}
