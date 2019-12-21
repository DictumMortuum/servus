package calendar

import (
	"bytes"
	"fmt"
	"github.com/twinj/uuid"
	"time"
)

type Vevent struct {
	Dtstart string
	Dtend   string
	Summary string
}

func f(t time.Time) string {
	return t.Format("20060102T150405")
}

func dtstamp() string {
	now := time.Now()
	return f(now.UTC())
}

func dtstart(shift, day int) string {
	year, month, _ := time.Now().UTC().Date()
	return f(time.Date(year, month+1, day+1, shift, 0, 0, 0, time.UTC))
}

func dtend(shift, day int) string {
	year, month, _ := time.Now().UTC().Date()
	return f(time.Date(year, month+1, day+1, shift+8, 0, 0, 0, time.UTC))
}

func (v Vevent) String() string {
	buf := bytes.NewBufferString("")

	fmt.Fprintf(buf, "BEGIN:VEVENT\n")
	fmt.Fprintf(buf, "UID:%s\n", uuid.NewV4())
	fmt.Fprintf(buf, "SEQUENCE:0\n")
	fmt.Fprintf(buf, "DTSTAMP:%s\n", dtstamp())
	fmt.Fprintf(buf, "DTSTART:%s\n", v.Dtstart)
	fmt.Fprintf(buf, "DTEND:%s\n", v.Dtend)
	fmt.Fprintf(buf, "SUMMARY:%s\n", v.Summary)
	fmt.Fprintf(buf, "END:VEVENT\n")

	return buf.String()
}
