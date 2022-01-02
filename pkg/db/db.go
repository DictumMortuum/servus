package db

import (
	"database/sql"
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/models"
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

func Insert(db *sqlx.DB, data models.Insertable) (*models.JsonNullInt64, error) {
	rs, err := db.NamedExec(data.Insert(), data)
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.JsonNullInt64{
		Int64: id,
		Valid: true,
	}, nil
}

func Exists(db *sqlx.DB, data models.Insertable) (*models.JsonNullInt64, error) {
	var id models.JsonNullInt64

	stmt, err := db.PrepareNamed(data.Exists())
	if err != nil {
		return nil, err
	}

	err = stmt.Get(&id, data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func InsertIfNotExists(db *sqlx.DB, data models.Insertable) (*models.JsonNullInt64, error) {
	id, err := Exists(db, data)
	if err != nil {
		return nil, err
	}

	if id == nil {
		id, err = Insert(db, data)
		if err != nil {
			return nil, err
		}

		return id, nil
	} else {
		return nil, nil
	}
}
