package db

import (
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DatabaseTypeConnect(t string, database string) (*sqlx.DB, error) {
	uri := config.App.Databases[database]
	db, err := sqlx.Connect(t, uri)
	if err != nil {
		return nil, err
	}

	return db, nil
}
