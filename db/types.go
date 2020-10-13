package db

import (
	"time"
)

type CalendarRow struct {
	Id           int `db:"calendar_id"`
	Index        int
	Raw          string
	DayName      string
	Uuid         string    `db:"uuid"`
	Date         time.Time `db:"date"`
	Shift        int       `db:"shift"`
	Summary      string    `db:"summary"`
	CreationDate time.Time `db:"cr_date"`
	Seq          int       `db:"sequence"`
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
