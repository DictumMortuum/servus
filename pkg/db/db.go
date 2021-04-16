package db

import (
	"github.com/DictumMortuum/servus/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Conn() (*sqlx.DB, error) {
	return DatabaseConnect("servus")
}

func DatabaseConnect(database string) (*sqlx.DB, error) {
	url := config.App.GetMariaDBConnection(database)
	db, err := sqlx.Connect("mysql", url)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RsIsEmpty(err error) bool {
	return err.Error() == "sql: no rows in result set"
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

	err := db.Select(&retval, "select * from tcalendar where date >= NOW() - interval 10 day order by date")
	if err != nil {
		return nil, err
	}

	return retval, nil
}

func GetShiftsWithoutNextcloudEntry(db *sqlx.DB) ([]CalendarRow, error) {
	retval := []CalendarRow{}

	sql := `select
		t.*
	from
		tcalendar t
	left join
		nextcloud.oc_calendarobjects n on t.uuid = n.uid
	where
		n.uid is null
	order by date
	limit 5
	`

	err := db.Select(&retval, sql)
	if err != nil {
		return nil, err
	}

	return retval, nil
}

func EventExists(db *sqlx.DB, day CalendarRow) (bool, error) {
	rows, err := db.NamedQuery(`select 1 from tcalendar where date=:date and shift=:shift`, day)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func CreateEvent(db *sqlx.DB, day CalendarRow) error {
	sql := `
	insert into tcalendar (
		uuid,
		date,
		shift,
		summary,
		description,
		cr_date,
		sequence
	) values (
		UUID(),
		:date,
		:shift,
		:summary,
		:description,
		NOW(),
		0
	) on duplicate key update
		shift=:shift,
		summary=:summary,
		description=:description,
		cr_date=NOW(),
		sequence=sequence+1,
		updated=1`

	_, err := db.NamedExec(sql, &day)
	if err != nil {
		return err
	}

	return nil
}
