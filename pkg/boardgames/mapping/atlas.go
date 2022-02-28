package mapping

import (
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/boardgames/atlas"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
)

func MapAtlas(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []mapping{}

	price, err := boardgames.GetPriceById(db, args.Id)
	if err != nil {
		return nil, err
	}

	atlas_results, err := atlas.Search(price.Name)
	if err != nil {
		return nil, err
	}

	for _, result := range atlas_results {
		name := transformName(result.Name)
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
	}

	return rs, nil
}

func SearchAtlasTerm(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []mapping{}

	atlas_results, err := atlas.Search(args.Data["term"].(string))
	if err != nil {
		return nil, err
	}

	for _, result := range atlas_results {
		name := transformName(result.Name)
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
	}

	return rs, nil
}
