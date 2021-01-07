package gas

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type FuelRow struct {
	Id           int            `db:"fuel_id"`
	Date         time.Time      `db:"date" form:"date" binding:"required"`
	CostPerLitre float64        `db:"cost_per_litre" form:"cost_per_litre" binding:"required"`
	Litre        float64        `db:"litre" form:"litre" binding:"required"`
	Cost         float64        `db:"cost" form:"cost" binding:"required"`
	Location     sql.NullString `db:"location" form:"location" binding:"required"`
}

type FuelStatsRow struct {
	Id                int           `db:"fuel_stats_id"`
	FuelId            int           `db:"fuel_id"`
	Kilometers        float32       `db:"km" form:"km" binding:"required"`
	LitreAverage      float32       `db:"litre_average" form:"litre_average" binding:"required"`
	Duration          time.Duration `db:"duration" form:"duration" binding:"required"`
	KilometersPerHour int           `db:"kmh" form:"kmh" binding:"required"`
}

type FuelJoinRow struct {
	Fuel      FuelRow      `db:"tfuel"`
	FuelStats FuelStatsRow `db:"tfuelstats"`
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
