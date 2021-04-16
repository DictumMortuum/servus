package generate

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"os"
	"text/template"
)

func Handler(c *gin.Context) {
	database, err := db.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	calendar, err := db.GetShifts(database)
	if err != nil {
		util.Error(c, err)
		return
	}

	t := template.Must(template.New("ics").Parse(`BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//dictummortuum.com//servus//EN{{ range . }}
BEGIN:VEVENT
UID:{{ .Uuid }}
SEQUENCE:{{ .Seq }}
DTSTAMP:{{ .Dtstamp }}
DTSTART:{{ .Dtstart }}
DTEND:{{ .Dtend }}
SUMMARY:{{ .Summary }}
DESCRIPTION:{{ .Description }}
END:VEVENT{{end}}
END:VCALENDAR
`))

	f, err := os.Create("/tmp/cal.ics")
	if err != nil {
		util.Error(c, err)
		return
	}

	t.Execute(f, calendar)
	c.FileAttachment("/tmp/cal.ics", "cal.ics")
}
