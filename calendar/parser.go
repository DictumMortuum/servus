package calendar

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"regexp"
	"strconv"
)

func Handler(c *gin.Context) {
	file, _ := c.FormFile("files")
	c.SaveUploadedFile(file, "/tmp/cal")
	data := Get("/tmp/cal")
	buffer := []byte(data.String())
	ioutil.WriteFile("/tmp/cal.ics", buffer, 0600)
	c.FileAttachment("/tmp/cal.ics", "cal.ics")
}

func checkCellForName(s string) bool {
	re1 := regexp.MustCompile("Ω")
	re2 := regexp.MustCompile("Φ")
	re3 := regexp.MustCompile("Σ")
	re4 := regexp.MustCompile("ΤΑΜ")

	if re1.FindStringSubmatch(s) != nil {
		return true
	}

	if re2.FindStringSubmatch(s) != nil {
		return true
	}

	if re3.FindStringSubmatch(s) != nil {
		return true
	}

	if re4.FindStringSubmatch(s) != nil {
		return true
	}

	return false
}

func findHeader(file *xlsx.File) (int, int, int) {
	for i, sheet := range file.Sheets {
		for j, row := range sheet.Rows {
			for k, cell := range row.Cells {
				s := cell.String()

				if s == "ΦΥΤΡΟΥ" || s == "ΦΥTΡΟΥ" || s == "3673" {
					return i, j, k
				}
			}
		}
	}

	return -1, -1, -1
}

func filterRow(row *xlsx.Row) []string {
	ret := []string{}

	for _, cell := range row.Cells {
		s := cell.String()

		if s == "3673" || checkCellForName(s) || s == "" {
			continue
		}

		ret = append(ret, s)
	}

	return ret
}

func transform(row []string) []string {
	ret := []string{}
	re := regexp.MustCompile("[MHΜΗ]{2}([0-9][0-9])")

	for _, cell := range row {
		s := re.FindStringSubmatch(cell)

		if s == nil {
			if !checkCellForName(cell) {
				ret = append(ret, "ΡΕΠΟ")
			}
		} else {
			i, _ := strconv.Atoi(s[1])
			ret = append(ret, strconv.Itoa(i))
		}
	}

	return ret
}

func Get(filename string) Vcalendar {
	xlFile, _ := xlsx.OpenFile(filename)
	i, j, _ := findHeader(xlFile)
	target := xlFile.Sheets[i].Rows[j]
	calendar := filterRow(target)

	events := []Vevent{}

	for day, shift := range transform(calendar) {
		var v Vevent
		s, _ := strconv.Atoi(shift)

		if shift == "ΡΕΠΟ" {
			v = Vevent{
				dtstart(s, day),
				dtstart(s, day),
				"Ρεπό",
			}
		} else {
			v = Vevent{
				dtstart(s, day),
				dtend(s, day),
				fmt.Sprintf("Βάρδια %s", shift),
			}
		}

		events = append(events, v)
	}

	return Vcalendar{
		events,
	}
}
