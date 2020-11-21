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

	err := db.Select(&retval, "select * from tcalendar where date >= NOW() - interval 0 day order by date")
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
		sequence=sequence+1`

	_, err := db.NamedExec(sql, &day)
	if err != nil {
		return err
	}

	return nil
}

func GetFuelIdsWithoutStats(db *sqlx.DB) ([]FuelRow, error) {
	retval := []FuelRow{}

	err := db.Select(&retval, "select * from tfuel where fuel_id not in (select fuel_id from tfuelstats) order by fuel_id")
	if err != nil {
		return nil, err
	}

	return retval, nil
}

func CreateFuelStats(db *sqlx.DB, row FuelStatsRow) error {
	ids, err := GetFuelIdsWithoutStats(db)
	if err != nil {
		return err
	}

	row.FuelId = ids[0].Id

	sql := `
	insert into tfuelstats (
		fuel_id,
		km,
		litre_average,
		duration,
		kmh
	) values (
		:fuel_id,
		:km,
		:litre_average,
		:duration,
		:kmh
	)`

	_, err = db.NamedExec(sql, &row)
	if err != nil {
		return err
	}

	return nil
}

func CreateFuel(db *sqlx.DB, row FuelRow) error {
	sql := `
	insert into tfuelstats (
		date,
		cost_per_litre,
		litre,
		cost,
		location
	) values (
		:date,
		:cost_per_litre,
		:litre,
		:cost,
		:location
	)`

	_, err := db.NamedExec(sql, &row)
	if err != nil {
		return err
	}

	return nil
}

func GetGas(db *sqlx.DB) ([]FuelJoinRow, error) {
	retval := []FuelJoinRow{}

	sql := `select
    f.fuel_id "tfuel.fuel_id",
    f.date "tfuel.date",
    f.cost_per_litre "tfuel.cost_per_litre",
    f.litre "tfuel.litre",
    f.cost "tfuel.cost",
    f.location "tfuel.location",
    s.fuel_stats_id "tfuelstats.fuel_stats_id",
    s.fuel_id "tfuelstats.fuel_id",
    s.km "tfuelstats.km",
    s.litre_average "tfuelstats.litre_average",
    s.duration "tfuelstats.duration",
    s.kmh "tfuelstats.kmh"
  from
    tfuel as f,
    tfuelstats s
  where
    f.fuel_id = s.fuel_id
	`

	err := db.Select(&retval, sql)
	if err != nil {
		return nil, err
	}

	return retval, nil
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
