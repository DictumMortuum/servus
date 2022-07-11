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
	detected := 0

	rconn, ch, q, err := rabbitmq.SetupQueue("prices-delay")
	if err != nil {
		return nil, err
	}
	defer rconn.Close()
	defer ch.Close()

	rows, err := updateBatch(db, store_id)
	if err != nil {
		return nil, err
	}

	log.Printf("Scraper %d started - resetting %d rows\n", store_id, rows)

	//https://boardgamegeek.com/xmlapi2/geeklist/125657
	rs, err := bgg.Geeklist(125657)
	if err != nil {
		return nil, err
	}

	for _, item := range rs.Items {
		urls := url.FindAllStringSubmatch(item.Body, -1)

		if len(urls) == 1 && item.ObjectId != 23953 {
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

			detected++
			err = rabbitmq.InsertQueueItem(ch, q, item)
			if err != nil {
				return nil, err
			}
		}
	}

	return map[string]interface{}{
		"name":     "Greek Guild",
		"id":       store_id,
		"scraped":  detected,
		"resetted": rows,
	}, nil
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
