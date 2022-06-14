package mapping

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/boardgames/atlas"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/lithammer/fuzzysearch/fuzzy"
	// google "github.com/rocketlaunchr/google-search"
	"sort"
	// "strings"
)

func updatePrice(db *sqlx.DB, payload models.Price) (bool, error) {
	sql := `
		update
			tboardgameprices
		set
			boardgame_id = :boardgame_id,
			levenshtein = :levenshtein,
			hamming = :hamming,
			batch = 1
		where id = :id
	`

	rs, err := db.NamedExec(sql, payload)
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func updatePriceNotMapped(db *sqlx.DB, payload models.Price) (bool, error) {
	sql := `
		update
			tboardgameprices
		set
			batch = 1
		where id = :id
	`

	rs, err := db.NamedExec(sql, payload)
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func getBoardgameNames(db *sqlx.DB) ([]string, error) {
	var rs []string

	sql := `
		select
			name
		from
			tboardgames
	`

	err := db.Select(&rs, sql)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func boardgameNameToId(db *sqlx.DB, name string) (*models.Boardgame, error) {
	var rs models.Boardgame

	sql := `
		select * from tboardgames where name = ?
	`

	err := db.QueryRowx(sql, name).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func getPricesWithoutMappings(db *sqlx.DB) ([]models.Price, error) {
	var rs []models.Price

	sql := `
		select
			p.*,
			NULL as rank,
			"" as thumb,
			"" as boardgame_name,
			s.name as store_name
		from
			tboardgameprices p,
			tboardgamestores s
		where
			p.boardgame_id is NULL and
			p.store_id = s.id
	`

	err := db.Select(&rs, sql)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getPrice(db *sqlx.DB, id int64) (*models.Price, error) {
	var rs models.Price

	err := db.QueryRowx(`
		select
			p.*,
			g.rank,
			IFNULL(g.thumb,"") as thumb,
			g.name as boardgame_name,
			s.name as store_name
		from
			tboardgameprices p,
			tboardgames g,
			tboardgamestores s
		where
			p.id = ? and
			p.boardgame_id = g.id and
			p.store_id = s.id
	`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func fuzzyFind(col []string) func(string) fuzzy.Ranks {
	return func(s string) fuzzy.Ranks {
		rs := fuzzy.RankFindNormalizedFold(s, col)
		sort.Sort(rs)
		l := len(rs)

		hi := 5
		if l < 5 {
			hi = l
		}

		return rs[0:hi]
	}
}

// func googlesearch(s string) ([]string, error) {
// 	retval := []string{}

// 	rs, err := google.Search(nil, s)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, item := range rs {
// 		if strings.Contains(item.URL, "boardgamegeek") {
// 			retval = append(retval, item.URL)
// 		}
// 	}

// 	return retval, nil
// }

func Map(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	price, err := getPrice(db, args.Id)
	if err != nil {
		return nil, err
	}

	name := boardgames.TransformName(price.Name)

	boardgames, err := getBoardgameNames(db)
	if err != nil {
		return nil, err
	}

	fn := fuzzyFind(boardgames)
	ranks := fn(name)

	if len(ranks) > 0 {
		// boardgame, err := boardgameNameToId(db, ranks[0].Target)
		// if err != nil {
		// 	return nil, err
		// }

		return ranks, nil
	}

	bgg_results, err := bgg.Search(name)
	if err != nil {
		return nil, err
	}

	if len(bgg_results) > 0 {
		return bgg_results, nil
	}

	atlas_results, err := atlas.Search(name)
	if err != nil {
		return nil, err
	}

	if len(atlas_results) > 0 {
		return atlas_results, nil
	}

	return atlas_results, nil

	// google_results, err := googlesearch(name)
	// if err != nil {
	// 	return nil, err
	// }

	// return google_results, nil
}

func MapAll(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	prices, err := getPricesWithoutMappings(db)
	if err != nil {
		return nil, err
	}

	for _, price := range prices {
		updatePriceNotMapped(db, price)
	}

	games, err := getBoardgameNames(db)
	if err != nil {
		return nil, err
	}

	fn := fuzzyFind(games)

	retval := []models.Price{}

	l := len(prices)

	for i, price := range prices {
		name := boardgames.TransformName(price.Name)
		ranks := fn(name)

		if len(ranks) > 0 {
			boardgame, err := boardgameNameToId(db, ranks[0].Target)
			if err != nil {
				return nil, err
			}

			fmt.Printf("[%5v/%v] %v to %v\n", i, l, ranks[0], boardgame.Id)

			price.BoardgameId = models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			}

			price.Hamming = Hamming(price.Name, ranks[0].Target)
			price.Levenshtein = ranks[0].Distance

			updatePrice(db, price)

			retval = append(retval, price)
		} else {
			fmt.Printf("[%5v/%v] %v not mapped\n", i, l, price.Name)

			updatePriceNotMapped(db, price)
		}
	}

	return retval, nil
}
