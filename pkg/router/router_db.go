package router

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type RouterRow struct {
	Id           int64     `db:"id"`
	Uptime       int64     `db:"uptime"`
	Date         time.Time `db:"date"`
	CrDate       time.Time `db:"cr_date"`
	MaxUp        int       `db:"max_up"`
	MaxDown      int       `db:"max_down"`
	CurrentUp    int       `db:"current_up"`
	CurrentDown  int       `db:"current_down"`
	InitialUp    int       `db:"initial_up"`
	InitialDown  int       `db:"initial_down"`
	CRCUp        int       `db:"crc_up"`
	CRCDown      int       `db:"crc_down"`
	FECUp        int       `db:"fec_up"`
	FECDown      int       `db:"fec_down"`
	SNRUp        int64     `db:"snr_up"`
	SNRDown      int64     `db:"snr_down"`
	DataUp       int64     `db:"data_up"`
	DataDown     int64     `db:"data_down"`
	Disconnected bool      `db:"disconnected"`
}

func CreateRouter(db *sqlx.DB, data RouterRow) error {
	sql := `
	insert into trouter (
		uptime,
		date,
		cr_date,
		max_up,
		max_down,
		current_up,
		current_down,
		initial_up,
		initial_down,
		crc_up,
		crc_down,
		fec_up,
		fec_down,
		snr_up,
		snr_down,
		data_up,
		data_down,
		disconnected
	) values (
		:uptime,
		:date,
		NOW(),
		:max_up,
		:max_down,
		:current_up,
		:current_down,
		:initial_up,
		:initial_down,
		:crc_up,
		:crc_down,
		:fec_up,
		:fec_down,
		:snr_up,
		:snr_down,
		:data_up,
		:data_down,
		:disconnected
	)`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRouter(db *sqlx.DB, data RouterRow) error {
	sql := `
	update trouter set
		uptime = :uptime,
		cr_date = NOW(),
		max_up = :max_up,
		max_down = :max_down,
		current_up = :current_up,
		current_down = :current_down,
		initial_up = IF(initial_up = 0, :current_up, :initial_up),
		initial_down = IF(initial_down = 0, :current_down, :initial_down),
		crc_up = :crc_up,
		crc_down = :crc_down,
		fec_up = :fec_up,
		fec_down = :fec_down,
		snr_up = :snr_up,
		snr_down = :snr_down,
		data_up = :data_up,
		data_down = :data_down,
		disconnected = :disconnected
	where id = :id`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func RouterExists(db *sqlx.DB, row RouterRow) (int64, error) {
	var id sql.NullInt64

	sql := `
	select
		max(id)
	from
		trouter
	where
		date >= :date - interval 1 minute and
		date <= :date + interval 1 minute
	`

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return -1, err
	}

	err = stmt.Get(&id, row)
	if err != nil {
		return -1, err
	}

	retval, err := id.Value()
	if err != nil {
		return -1, err
	}

	if retval == nil {
		return -1, err
	}

	return retval.(int64), nil
}

func getLatestRouter(db *sqlx.DB) (*RouterRow, error) {
	var retval RouterRow

	err := db.QueryRowx(`
	select
		*
	from
		trouter
	order by
		cr_date desc
	limit 1
	`).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}

func getRouter(db *sqlx.DB, id int64) (*RouterRow, error) {
	var retval RouterRow

	err := db.QueryRowx(`
	select
		*
	from
		trouter
	where
		id = ?
	`, id).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}
