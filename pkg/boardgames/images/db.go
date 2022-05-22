package images

import (
	"github.com/jmoiron/sqlx"
)

func getBoardgameUrls(db *sqlx.DB, id int64) ([]string, error) {
	var rs []string

	sql := `
		select
			thumb
		from
			tboardgames
		where
			id = ?
			and thumb is not null
	`

	err := db.Select(&rs, sql, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getPriceUrls(db *sqlx.DB, id int64) ([]string, error) {
	var rs []string

	sql := `
		select
			store_thumb
		from
			tboardgameprices
		where
			boardgame_id = ?
	`

	err := db.Select(&rs, sql, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
