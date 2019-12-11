package main

import (
	"fmt"
	"github.com/DictumMortuum/servus/calendar"
	"github.com/gin-gonic/gin"
)

func startpageHandler(c *gin.Context) {
	data, _ := Asset("html/index.html")
	c.Writer.Write(data)
}

func calendarHandler(c *gin.Context) {
	file, _ := c.FormFile("file")
	c.SaveUploadedFile(file, "/tmp/cal")
	tmp := calendar.Get("/tmp/cal")
	c.String(200, fmt.Sprintf("%v", tmp))
}

func main() {
	r := gin.New()
	r.StaticFS("/assets", assetFS())
	r.GET("/startpage", startpageHandler)
	r.POST("/calendar", calendarHandler)
	r.Run("127.0.0.1:1234")
}
