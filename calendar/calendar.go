package calendar

import (
	"bytes"
	"fmt"
)

type Vcalendar struct {
	Events []Vevent
}

func (v Vcalendar) String() string {
	buf := bytes.NewBufferString("")

	fmt.Fprintf(buf, "BEGIN:VCALENDAR\n")
	fmt.Fprintf(buf, "VERSION:2.0\n")
	fmt.Fprintf(buf, "PRODID:-//dictummortuum.com//servus//EN\n")

	for _, event := range v.Events {
		fmt.Fprintf(buf, "%s", event)
	}

	fmt.Fprintf(buf, "END:VCALENDAR\n")

	return buf.String()
}
