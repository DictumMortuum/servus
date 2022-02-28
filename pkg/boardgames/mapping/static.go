package mapping

import (
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
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
