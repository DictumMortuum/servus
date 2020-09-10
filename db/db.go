package db

import (
	"database/sql"
	"github.com/DictumMortuum/servus/config"
	_ "github.com/go-sql-driver/mysql"
)

func Conn() (*sql.DB, error) {
	user := config.App.Database.Username
	pass := config.App.Database.Password
	db, err := sql.Open("mysql", user+":"+pass+"@tcp(127.0.0.1:3306)/servus?parseTime=true")

	if err != nil {
		return nil, err
	}

	return db, nil
}
