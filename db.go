package db

import (
	"bytes"
	"fmt"
	"github.com/DictumMortuum/servus/config"
	"github.com/DictumMortuum/servus/generic"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"text/template"
)

func Conn() (*sqlx.DB, error) {
	url := config.App.GetMariaDBConnection("servus?parseTime=true")
	db, err := sqlx.Connect("mysql", url)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetShiftList(c *gin.Context) {
	err, args := generic.ParseArgs(c)
	if err != nil {
		generic.Fail(c, err)
		return
	}

	db, err := Conn()
	if err != nil {
		generic.Fail(c, err)
		return
	}
	defer db.Close()

	log.Println(args)

	var total int
	err = db.Get(&total, "select count(*) from tcalendar where date >= NOW() - interval 1 day")
	if err != nil {
		generic.Fail(c, err)
		return
	}

	sql := `
select
	*
from
	tcalendar
where date >= NOW() - interval 1 day
{{ if eq (len .Sort) 2 }}
order by {{ index .Sort 0 }} {{ index .Sort 1 }}
{{ end }}
{{ if eq (len .Range) 2  }}
limit {{ index .Range 0 }}, {{ .Page }}
{{ end }}
`

	var tpl bytes.Buffer
	t := template.Must(template.New("shiftlist").Parse(sql))
	err = t.Execute(&tpl, args)
	if err != nil {
		generic.Fail(c, err)
		return
	}

	data := []CalendarRow{}
	err = db.Select(&data, tpl.String())
	if err != nil {
		generic.Fail(c, err)
		return
	}

	c.Header("Content-Range", fmt.Sprintf("%d-%d/%d", args.Range[0], args.Range[1], total))
	c.JSON(http.StatusOK, data)
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

func LinkExists(db *sqlx.DB, row LinkRow) (bool, error) {
	rows, err := db.NamedQuery(`select 1 from tlink where url=:url and host=:host and user=:user`, row)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}
