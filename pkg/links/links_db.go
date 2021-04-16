package links

import (
	"github.com/jmoiron/sqlx"
)

type LinkRow struct {
	Url  string `db:"url" form:"url" binding:"required"`
	Host string `db:"host" form:"host" binding:"required"`
	User string `db:"user" form:"user" binding:"required"`
}

func CreateLink(db *sqlx.DB, row LinkRow) error {
	sql := `
	insert into tlink (
		url,
		host,
		user
	) values (
		:url,
		:host,
		:user
	)`

	_, err := db.NamedExec(sql, &row)
	if err != nil {
		return err
	}

	return nil
}

func LinkExists(db *sqlx.DB, row LinkRow) (bool, error) {
	rows, err := db.NamedQuery(`select 1 from tlink where url=:url and host=:host and user=:user`, row)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}
