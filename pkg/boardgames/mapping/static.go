package mapping

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

func getBoardgameName(db *sqlx.DB, name string) (*mapping, error) {
	var rs mapping

	q := `
		select * from tboardgamepricesmap where name = ?
	`

	err := db.QueryRowx(q, name).StructScan(&rs)
	// if err == sql.ErrNoRows {
	// 	return nil, nil
	// } else if err != nil {
	// 	return nil, err
	// }
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func MapStatic(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	price, err := boardgames.GetPriceById(db, args.Id)
	if err != nil {
		return nil, err
	}

	match, err := getBoardgameName(db, price.Name)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func MapAllStatic(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	prices, err := getPricesWithoutMappings(db)
	if err != nil {
		return nil, err
	}

	ignored := []string{
		"dobble",
		"gravitrax",
		"story cubes",
		"puzzle",
		"monopoly",
		"desyllas",
		"δεσύλλας",
		"δεσυλλας",
		"σπαζοκεφαλιά",
		"σπαζοκεφαλια",
		"think fun",
		"sleeves",
		"υπερατού",
		"sudoku",
		"thinkfun",
		"zito!",
		"top trumps",
		"funkoverse",
		"πλαστικοποιημένη",
		"κουτί για κάρτες",
		"προστατευτικά καρτών",
		"pokemon tcg",
		"yu-gi-oh",
		"κουτί για κάρτες",
		"similo",
		"magic the gathering",
		"κουτί για χαρτιά",
		"orchard toys",
		"desyllas",
	}

	retval := []models.Price{}
	l := len(prices)

	for _, price := range prices {
		tmp := strings.ToLower(price.Name)

		for _, ignore := range ignored {
			if strings.Contains(tmp, ignore) {
				price.BoardgameId = models.JsonNullInt64{
					Int64: 23953,
					Valid: true,
				}

				updatePrice(db, price)
				break
			}
		}
	}

	for i, price := range prices {
		match, _ := getBoardgameName(db, boardgames.TransformName(price.Name))
		if match != nil {
			fmt.Printf("[%5v/%v] %v to %v\n", i, l, price.Name, match.BoardgameId)

			price.BoardgameId = models.JsonNullInt64{
				Int64: match.BoardgameId,
				Valid: true,
			}

			updatePrice(db, price)

			retval = append(retval, price)
		}
	}

	return retval, nil
}
