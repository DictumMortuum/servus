package router

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type RouterRow struct {
	Uptime      int64     `db:"uptime"`
	Date        time.Time `db:"date"`
	MaxUp       int       `db:"max_up"`
	MaxDown     int       `db:"max_down"`
	CurrentUp   int       `db:"current_up"`
	CurrentDown int       `db:"current_down"`
	CRCUp       int       `db:"crc_up"`
	CRCDown     int       `db:"crc_down"`
	FECUp       int       `db:"fec_up"`
	FECDown     int       `db:"fec_down"`
	SNRUp       int64     `db:"snr_up"`
	SNRDown     int64     `db:"snr_down"`
	DataUp      int64     `db:"data_up"`
	DataDown    int64     `db:"data_down"`
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
		crc_up,
		crc_down,
		fec_up,
		fec_down,
		snr_up,
		snr_down,
		data_up,
		data_down
	) values (
		:uptime,
		:date,
		NOW(),
		:max_up,
		:max_down,
		:current_up,
		:current_down,
		:crc_up,
		:crc_down,
		:fec_up,
		:fec_down,
		:snr_up,
		:snr_down,
		:data_up,
		:data_down
	)`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func RouterExists(db *sqlx.DB, row RouterRow) (bool, error) {
	rows, err := db.NamedQuery(`select 1 from trouter where date=:date`, row)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}
