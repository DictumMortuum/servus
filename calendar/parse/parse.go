package parse

import (
	"github.com/DictumMortuum/servus/calendar"
	"github.com/DictumMortuum/servus/db"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type formData struct {
	Year  string                `form:"year" binding:"required"`
	Month string                `form:"month" binding:"required"`
	File  *multipart.FileHeader `form:"files" binding:"required"`
}

func parseXlsx(form formData) error {
	database, err := db.Conn()
	if err != nil {
		return err
	}
	defer database.Close()

	sheet := calendar.Sheet{}
	err = sheet.New("/tmp/" + form.File.Filename)
	if err != nil {
		return err
	}

	m, err := strconv.Atoi(form.Month)
	if err != nil {
		return err
	}

	y, err := strconv.Atoi(form.Year)
	if err != nil {
		return err
	}

	for _, person := range sheet.People {
		if person.Id == "3673" {
			for _, day := range person.Calendar {
				coworkers := calendar.GetCoworkers(day, sheet.Header)
				day.Date = day.GetDate(y, m)
				day.Shift = calendar.RawToShift(day.Raw)
				day.Summary = calendar.FormatShift(day.Shift)
				day.Description = calendar.FormatCoworkers(coworkers)

				err := db.CreateEvent(database, day)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func Handler(c *gin.Context) {
	var form formData

	err := c.ShouldBind(&form)
	if err != nil {
		c.SetCookie("error", err.Error(), 3600, "", "", false, true)
		c.Redirect(http.StatusMovedPermanently, "/calendar")
		return
	}

	c.SaveUploadedFile(form.File, "/tmp/"+form.File.Filename)

	err = parseXlsx(form)
	if err != nil {
		c.SetCookie("error", err.Error(), 3600, "", "", false, true)
		c.Redirect(http.StatusMovedPermanently, "/calendar")
		return
	}

	state := gin.H{
		"title": "Calendar",
		"primary": map[string]interface{}{
			"enabled": true,
			"desc":    "Upload xlsx",
		},
		"secondary": map[string]interface{}{
			"enabled": false,
		},
		"redirect": true,
	}

	//c.Redirect(http.StatusMovedPermanently, "/calendar")
	c.HTML(http.StatusOK, "parse_successful.html", state)
}

func Render(c *gin.Context) {
	cookie, _ := c.Cookie("error")
	c.SetCookie("error", "", 3600, "", "", false, true)

	state := gin.H{
		"title": "Calendar",
		"primary": map[string]interface{}{
			"enabled": true,
			"desc":    "Upload xlsx",
		},
		"secondary": map[string]interface{}{
			"enabled": false,
		},
		"error": cookie,
	}

	currentYear, currentMonth, _ := time.Now().Date()
	months := []map[string]interface{}{}
	years := []map[string]interface{}{}

	for i := 1; i <= 12; i++ {
		months = append(months, map[string]interface{}{
			"Index":    i,
			"Name":     time.Month(i).String(),
			"Selected": int(currentMonth)+1 == i,
		})
	}

	for i := -5; i < 5; i++ {
		years = append(years, map[string]interface{}{
			"Index":    currentYear + i,
			"Selected": i == 0,
		})
	}

	state["months"] = months
	state["years"] = years

	database, err := db.Conn()
	if err != nil {
		state["error"] = err.Error()
		c.HTML(http.StatusOK, "parse.html", state)
		return
	}
	defer database.Close()

	calendar, err := db.GetFutureShifts(database)
	if err != nil {
		state["error"] = err.Error()
		c.HTML(http.StatusOK, "parse.html", state)
		return
	}

	state["calendar"] = calendar

	c.HTML(http.StatusOK, "parse.html", state)
}
