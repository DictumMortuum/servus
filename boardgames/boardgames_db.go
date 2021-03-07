package boardgames

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type PlayersRow struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type PlaysRow struct {
	Id        int64     `db:"id" json:"id"`
	CrDate    time.Time `db:"cr_date" json:"cr_date"`
	Date      time.Time `db:"date" json:"date"`
	Boardgame string    `db:"boardgame"`
}

func createPlayer(db *sqlx.DB, data PlayersRow) error {
	sql := `
	insert into tboardgameplayers (
		name
	) values (
		:name
	)`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func getPlayer(db *sqlx.DB, id int64) (*PlayersRow, error) {
	var retval PlayersRow

	err := db.QueryRowx(`select * from tboardgameplayers where id = ?`, id).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}

func getPlayerByName(db *sqlx.DB, name string) (*PlayersRow, error) {
	var retval PlayersRow

	err := db.QueryRowx(`select * from tboardgameplayers where name = ?`, name).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
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
