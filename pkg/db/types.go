package db

import (
	"bytes"
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/tealeg/xlsx/v3"
	"net/http"
	"text/template"
	"time"
)

type CalendarRow struct {
	Id           int `db:"id"`
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
	Updated      bool      `db:"updated"`
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

func (c CalendarRow) ToICS() string {
	var output bytes.Buffer

	t := template.Must(template.New("ics").Parse(`BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//dictummortuum.com//servus//EN
BEGIN:VEVENT
UID:{{ .Uuid }}
SEQUENCE:{{ .Seq }}
DTSTART:{{ .Dtstart }}
DTEND:{{ .Dtend }}
SUMMARY:{{ .Summary }}
DESCRIPTION:{{ .Description }}
END:VEVENT
END:VCALENDAR
`))

	t.Execute(&output, c)
	return output.String()
}

func (c CalendarRow) ToCalDavServer() error {
	ics := []byte(c.ToICS())

	req, err := http.NewRequest("PUT", config.App.Calendar.Server+"/"+c.Uuid+".ics", bytes.NewBuffer(ics))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.App.Calendar.Username, config.App.Calendar.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
