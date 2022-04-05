package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
)

func ScrapeFantasyGate(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []models.Price{}

	collector := colly.NewCollector(
		colly.AllowedDomains("www.fantasygate.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".sblock4", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".jshop_price")

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    2,
			StoreThumb: e.ChildAttr(".jshop_img", "src"),
			Stock:      childHasClass(e, ".btn", "button_buy"),
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
		})
	})

	collector.Post("https://www.fantasygate.gr/strategygames", map[string]string{
		"limit": "99999",
	})

	collector.Post("https://www.fantasygate.gr/family-games", map[string]string{
		"limit": "99999",
	})

	collector.Post("https://www.fantasygate.gr/cardgames", map[string]string{
		"limit": "99999",
	})

	collector.Wait()

	err := updateBatch(db, 2)
	if err != nil {
		return nil, err
	}

	for _, item := range rs {
		item.BoardgameId = models.JsonNullInt64{
			Int64: -1,
			Valid: false,
		}

		id, err := findPrice(db, item)
		if err != nil {
			return nil, err
		}

		if id == nil {
			_, err := create(db, item)
			if err != nil {
				return nil, err
			}
		} else {
			item.Id = id.Int64
			_, err := update(db, item)
			if err != nil {
				return nil, err
			}
		}
	}

	return rs, nil
}
