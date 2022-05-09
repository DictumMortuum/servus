package mapping

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func getBoardgame(db *sqlx.DB, name string) (*mapping, error) {
	var rs mapping

	q := `
		select id, id as boardgame_id, name from tboardgames where name = ?
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

func SearchMaps(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	name := args.Query

	match, err := getBoardgame(db, name)
	if err != nil {
		return nil, err
	}

	return match, nil
}
