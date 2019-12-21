package main

import (
	"github.com/DictumMortuum/servus/calendar"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func staticPage(path string) func(*gin.Context) {
	return func(c *gin.Context) {
		data, _ := Asset(path)
		c.Writer.Write(data)
	}
}

func calendarHandler(c *gin.Context) {
	file, _ := c.FormFile("files")
	c.SaveUploadedFile(file, "/tmp/cal")
	data := calendar.Get("/tmp/cal")
	buffer := []byte(data.String())
	ioutil.WriteFile("/tmp/cal.ics", buffer, 0600)
	c.FileAttachment("/tmp/cal.ics", "cal.ics")
}

func main() {
	r := gin.New()
	r.StaticFS("/assets", assetFS())
	r.GET("/startpage", staticPage("html/index.html"))
	r.GET("/calendar", staticPage("html/calendar.html"))
	r.POST("/calendar/generate", calendarHandler)
	r.Run("127.0.0.1:1234")
}
