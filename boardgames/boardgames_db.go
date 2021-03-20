package boardgames

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type PlaysRow struct {
	Id        int64     `db:"id" json:"id"`
	CrDate    time.Time `db:"cr_date" json:"cr_date"`
	Date      time.Time `db:"date" json:"date"`
	Boardgame string    `db:"boardgame"`
}

func createDuelPlay(db *sqlx.DB, data PlaysRow) (int64, error) {
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
