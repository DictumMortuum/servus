package db

import (
	"database/sql"
	"github.com/tealeg/xlsx/v3"
	"time"
)

type FuelRow struct {
	Id           int            `db:"fuel_id"`
	Date         time.Time      `db:"date"`
	CostPerLitre float64        `db:"cost_per_litre"`
	Litre        float64        `db:"litre"`
	Cost         float64        `db:"cost"`
	Location     sql.NullString `db:"location"`
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

type CalendarRow struct {
	Id           int `db:"calendar_id"`
	Index        int
	Raw          string
	DayName      string
	Uuid         string    `db:"uuid"`
	Date         time.Time `db:"date"`
	Shift        int       `db:"shift"`
	Summary      string    `db:"summary"`
	Description  string    `db:"description"`
	CreationDate time.Time `db:"cr_date"`
	Seq          int       `db:"sequence"`
	Row          *xlsx.Row
	X            int
	Y            int
}

func (c *CalendarRow) Dtstamp() string {
	return c.CreationDate.Format("20060102T150405")
}

func (c CalendarRow) Dtstart() string {
	start := c.Date

	if c.Shift > 0 {
		start = start.Add(time.Hour * time.Duration(c.Shift))
	}

	return start.Format("20060102T150405")
}

func (c CalendarRow) Dtend() string {
	end := c.Date

	if c.Shift > 0 {
		end = end.Add(time.Hour * time.Duration(c.Shift+8))
	} else {
		end = end.Add(time.Hour*time.Duration(24) - time.Second)
	}

	return end.Format("20060102T150405")
}

func (d *CalendarRow) GetDate(year, month int) time.Time {
	return time.Date(year, time.Month(month), d.Index, 0, 0, 0, 0, time.UTC)
}
