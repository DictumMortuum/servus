package calendar

import (
	"database/sql"
	"fmt"
	"github.com/DictumMortuum/servus/db"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const create_event = `
insert into tcalendar (
	uuid,
	date,
	shift,
	summary,
	cr_date,
	sequence
) values (
	UUID(),
	?,
	?,
	?,
	NOW(),
	0
) on duplicate key update
	shift=?,
	summary=?,
	cr_date=NOW(),
	sequence=sequence+1`

func ParseHandler(c *gin.Context) {
	file, _ := c.FormFile("files")
	c.SaveUploadedFile(file, "/tmp/cal")

	db, err := db.Conn()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	for day, shift := range transform("/tmp/cal") {
		err = createEvent(db, c, day, shift)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
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

func transform(file string) []string {
	xlFile, _ := xlsx.OpenFile(file)
	i, j, _ := findHeader(xlFile)
	target := xlFile.Sheets[i].Rows[j]
	row := filterRow(target)
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

func createEvent(db *sql.DB, c *gin.Context, day int, shift string) error {
	stmt, err := db.Prepare(create_event)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var summary string

	if shift == "ΡΕΠΟ" {
		summary = "Ρεπό"
	} else {
		summary = fmt.Sprintf("Βάρδια %s", shift)
	}

	year := c.PostForm("year")
	month := c.PostForm("month")
	m, _ := strconv.Atoi(month)
	y, _ := strconv.Atoi(year)
	date := time.Date(y, time.Month(m), day+1, 0, 0, 0, 0, time.UTC)

	s, _ := strconv.Atoi(shift)

	_, err = stmt.Exec(date, s, summary, s, summary)
	if err != nil {
		return err
	}

	return nil
}
