package boardgames

import (
	"github.com/jmoiron/sqlx"
)

type PlayersRow struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

func CreatePlayer(db *sqlx.DB, data PlayersRow) error {
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

func GetPlayer(db *sqlx.DB, id int64) (*PlayersRow, error) {
	var retval PlayersRow

	err := db.QueryRowx(`select * from tboardgameplayers where id = ?`, id).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}

func GetPlayerByName(db *sqlx.DB, name string) (*PlayersRow, error) {
	var retval PlayersRow

	err := db.QueryRowx(`select * from tboardgameplayers where name = ?`, name).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}
