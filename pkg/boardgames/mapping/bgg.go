package mapping

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

func MapBGG(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []mapping{}

	price, err := boardgames.GetPriceById(db, args.Id)
	if err != nil {
		return nil, err
	}

	name := transformName(price.Name)
	log.Println(name)
	bgg_results, err := bgg.Search(name)
	if err != nil {
		return nil, err
	}

	for _, result := range bgg_results {
		rs = append(rs, mapping{
			Id:          -1,
			BoardgameId: result.Id,
			Name:        result.Name.Value,
		})
	}

	return rs, nil
}

func SearchBGGTerm(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []mapping{}

	name := transformName(args.Data["term"].(string))
	log.Println(name)
	bgg_results, err := bgg.Search(name)
	if err != nil {
		return nil, err
	}

	for _, result := range bgg_results {
		rs = append(rs, mapping{
			Id:          -1,
			BoardgameId: result.Id,
			Name:        result.Name.Value,
		})
	}

	return rs, nil
}

func MapAllBgg(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	prices, err := getPricesWithoutMappings(db)
	if err != nil {
		return nil, err
	}

	for _, price := range prices {
		updatePriceNotMapped(db, price)
	}

	retval := []models.Price{}

	l := len(prices)

	for i, price := range prices {
		name := transformName(price.Name)

		bgg_results, err := bgg.Search(name)
		if err != nil {
			return nil, err
		}

		if len(bgg_results) == 1 {
			price.BoardgameId = models.JsonNullInt64{
				Int64: bgg_results[0].Id,
				Valid: true,
			}

			fmt.Printf("[%5v/%v] %v to %v\n", i, l, price.Name, bgg_results[0].Id)

			updatePrice(db, price)

			retval = append(retval, price)
		}

		time.Sleep(3 * time.Second)
	}

	return retval, nil
}
