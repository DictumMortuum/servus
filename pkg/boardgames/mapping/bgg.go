package mapping

import (
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
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
