package db

import (
	"github.com/DictumMortuum/servus/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Conn() (*sqlx.DB, error) {
	return DatabaseConnect("mariadb")
}

func DatabaseConnect(database string) (*sqlx.DB, error) {
	uri := config.App.Databases[database]
	db, err := sqlx.Connect("mysql", uri)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RsIsEmpty(err error) bool {
	return err.Error() == "sql: no rows in result set"
}
