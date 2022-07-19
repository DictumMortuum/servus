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

	match, err := getBoardgameName(db, boardgames.TransformName(price.Name))
	if match != nil {
		price.BoardgameId = models.JsonNullInt64{
			Int64: match.BoardgameId,
			Valid: true,
		}

		updatePrice(db, *price)
	}

	return match, nil
}

func MapPrice(db *sqlx.DB, id int64) error {
	price, err := boardgames.GetPriceById(db, id)
	if err != nil {
		return err
	}

	match, err := getBoardgameName(db, boardgames.TransformName(price.Name))
	if match != nil {
		price.BoardgameId = models.JsonNullInt64{
			Int64: match.BoardgameId,
			Valid: true,
		}

		updatePrice(db, *price)
	}

	return nil
}

func Ignore(price models.Price) bool {
	tmp := strings.ToLower(price.Name)

	ignored := []string{
		"σετ ζαριών",
		"oργανωτής επιτραπέζιων",
		"dobble",
		"gravitrax",
		"story cubes",
		"puzzle",
		"monopoly",
		"desyllas",
		"δεσύλλας",
		"as company",
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
		"playing cards",
		"tcg",
		"jenga",
		"memory game",
		"tarot cards",
		"tarot deck",
		"chess set",
		"roleplaying",
		"claymates",
		"hasbro",
		"trivia quiz",
		"polly pocket",
		"ρόλων",
		"κουτί τράπουλας",
		"παζλ",
		"toploader",
		"τράπουλα",
		"gigamic",
		"50/50 games",
		"brainbox",
		"κουτί τράπουλας",
		"τράπουλα μονή",
		"ρόλων",
		"as επιτραπέζιο",
		"ultra pro",
		"as games",
	}

	for _, ignore := range ignored {
		if strings.Contains(tmp, ignore) {
			return true
		}
	}

	return false
}

func MapAllStatic(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	prices, err := getPricesWithoutMappings(db)
	if err != nil {
		return nil, err
	}

	retval := []models.Price{}
	l := len(prices)

	for _, price := range prices {
		if Ignore(price) {
			price.BoardgameId = models.JsonNullInt64{
				Int64: 23953,
				Valid: true,
			}

			updatePrice(db, price)
			break
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
