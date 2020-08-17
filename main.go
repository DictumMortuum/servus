package main

import (
	"github.com/DictumMortuum/servus/calendar"
	"github.com/DictumMortuum/servus/youtube"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func calendarHandler(c *gin.Context) {
	file, _ := c.FormFile("files")
	c.SaveUploadedFile(file, "/tmp/cal")
	data := calendar.Get("/tmp/cal")
	buffer := []byte(data.String())
	ioutil.WriteFile("/tmp/cal.ics", buffer, 0600)
	c.FileAttachment("/tmp/cal.ics", "cal.ics")
}

func main() {
	mode := os.Getenv("GIN_MODE")

	if mode == "release" {
		gin.DisableConsoleColor()
		f, _ := os.Create("/var/log/servus.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}

	err := os.MkdirAll("/var/lib/servus", 0755)
	if err != nil {
		log.Println(err)
	}

	r := gin.Default()
	r.POST("/calendar", calendarHandler)
	r.GET("/youtube/:url", youtube.Handler)
	r.Run("127.0.0.1:1234")
}
