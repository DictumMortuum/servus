package boardgames

import (
	"github.com/jmoiron/sqlx"
)

func createDuelPlay(db *sqlx.DB, data PlayModel) (int64, error) {
	res, err := db.NamedExec(`
	insert into tboardgameplays (
		cr_date,
		date
	) values (
		NOW(),
		:date,
		:boardgame
	)`, &data)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
