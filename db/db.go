package db

import (
	"github.com/DictumMortuum/servus/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Conn() (*sqlx.DB, error) {
	url := config.App.GetMariaDBConnection("servus?parseTime=true")
	db, err := sqlx.Connect("mysql", url)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetShifts(db *sqlx.DB) ([]CalendarRow, error) {
	retval := []CalendarRow{}

	err := db.Select(&retval, "select * from tcalendar order by date")
	if err != nil {
		return nil, err
	}

	return retval, nil
}

func GetFutureShifts(db *sqlx.DB) ([]CalendarRow, error) {
	retval := []CalendarRow{}

	err := db.Select(&retval, "select * from tcalendar where date > NOW() - interval 0 day order by date")
	if err != nil {
		return nil, err
	}

	return retval, nil
}

func CreateEvent(db *sqlx.DB, day CalendarRow) error {
	sql := `
	insert into tcalendar (
		uuid,
		date,
		shift,
		summary,
		cr_date,
		sequence
	) values (
		UUID(),
		:date,
		:shift,
		:summary,
		NOW(),
		0
	) on duplicate key update
		shift=:shift,
		summary=:summary,
		cr_date=NOW(),
		sequence=sequence+1`

	_, err := db.NamedExec(sql, &day)
	if err != nil {
		return err
	}

	return nil
}
