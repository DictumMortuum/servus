package calendar

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"regexp"
	"strconv"
)

func findHeader(file *xlsx.File) (int, int, int) {
	for i, sheet := range file.Sheets {
		for j, row := range sheet.Rows {
			for k, cell := range row.Cells {
				s := cell.String()

				if s == "ΦΥΤΡΟΥ" {
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

		if s == "3673" || s == "ΦΥΤΡΟΥ" || s == "ΘΕΩΝΗ" || s == "ΤΑΜΙΑΣ" || s == "" {
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
			ret = append(ret, "ΡΕΠΟ")
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
