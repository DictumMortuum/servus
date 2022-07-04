package search

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/jmoiron/sqlx"
	"log"
	"regexp"
)

var (
	url = regexp.MustCompile(`https://boardgamegeek.com/geekmarket/product/([0-9]+)`)
)

func ScrapeGreekGuild(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	store_id := int64(29)

	log.Printf("Scraper %d started\n", store_id)

	rconn, ch, q, err := rabbitmq.SetupQueue("prices-delay")
	if err != nil {
		return nil, err
	}
	defer rconn.Close()
	defer ch.Close()

	err = updateBatch(db, store_id)
	if err != nil {
		return nil, err
	}

	rs, err := bgg.Geeklist(125657)
	if err != nil {
		return nil, err
	}

	for _, item := range rs.Items {
		urls := url.FindAllStringSubmatch(item.Body, -1)

		if len(urls) == 1 {
			productId := urls[0][1]

			item := models.Price{
				Name:       item.Name,
				StoreId:    store_id,
				StoreThumb: "",
				Stock:      0,
				Url:        fmt.Sprintf("https://boardgamegeek.com/geeklist/125657/greek-guilds-games-sale?itemid=%d", item.ItemId),
				ExtraId:    item.ItemId,
				ProductId:  productId,
				BoardgameId: models.JsonNullInt64{
					Int64: item.ObjectId,
					Valid: true,
				},
			}

			err = rabbitmq.InsertQueueItem(ch, q, item)
			if err != nil {
				return nil, err
			}
		}
	}

	return rs, nil
}

// type GeeklistItem struct {
// 	XMLName xml.Name `xml:"item" json:"-"`
// 	Id      int64    `xml:"objectid,attr" json:"id"`
// 	ItemId  int64    `xml:"id,attr" json:"item_id"`
// 	Name    string   `xml:"objectname,attr" json:"name"`
// 	Date    string   `xml:"editdate,attr" json:"date"`
// 	Body    string   `xml:"body" json:"body"`
// }

// type GeeklistRs struct {
// 	XMLName xml.Name       `xml:"geeklist" json:"-"`
// 	Items   []GeeklistItem `xml:"item" json:"items"`
// }
