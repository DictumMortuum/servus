package boardgames

import (
	"github.com/jmoiron/sqlx"
)

type PlayersRow struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

func GetPlayer(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getPlayer(db, id)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func GetManyPlayers(db *sqlx.DB, ids []int64) (interface{}, int, error) {
	var rs []PlayersRow

	sql := `
	select
		* 
	from
		tboardgameplayers
	where
		id in (?)`

	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return nil, -1, err
	}

	err = db.Select(&rs, db.Rebind(query), args...)
	if err != nil {
		return nil, -1, err
	}

	return rs, len(rs), nil
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
