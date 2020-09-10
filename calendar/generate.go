package calendar

import (
	"bytes"
	"fmt"
	"github.com/DictumMortuum/servus/db"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func GenerateHandler(c *gin.Context) {
	db, err := db.Conn()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rs, err := db.Query("select uuid, date, shift, summary, cr_date, sequence from tcalendar order by date")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buf := bytes.NewBufferString("")

	fmt.Fprintf(buf, "BEGIN:VCALENDAR\n")
	fmt.Fprintf(buf, "VERSION:2.0\n")
	fmt.Fprintf(buf, "PRODID:-//dictummortuum.com//servus//EN\n")

	for rs.Next() {
		var uuid, summary string
		var shift, sequence int
		var date, cr_date time.Time

		err = rs.Scan(&uuid, &date, &shift, &summary, &cr_date, &sequence)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Fprintf(buf, "BEGIN:VEVENT\n")
		fmt.Fprintf(buf, "UID:%s\n", uuid)
		fmt.Fprintf(buf, "SEQUENCE:%d\n", sequence)
		fmt.Fprintf(buf, "DTSTAMP:%s\n", cr_date.Format("20060102T150405"))
		fmt.Fprintf(buf, "DTSTART:%s\n", date.Add(time.Hour*time.Duration(shift)).Format("20060102T150405"))
		fmt.Fprintf(buf, "DTEND:%s\n", date.Add(time.Hour*time.Duration(shift+8)).Format("20060102T150405"))
		fmt.Fprintf(buf, "SUMMARY:%s\n", summary)
		fmt.Fprintf(buf, "END:VEVENT\n")
	}

	fmt.Fprintf(buf, "END:VCALENDAR\n")

	ioutil.WriteFile("/tmp/cal.ics", []byte(buf.String()), 0600)
	c.FileAttachment("/tmp/cal.ics", "cal.ics")
}
