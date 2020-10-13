package main

import (
	"github.com/DictumMortuum/servus/calendar"
	"github.com/DictumMortuum/servus/calendar/generate"
	"github.com/DictumMortuum/servus/calendar/parse"
	"github.com/DictumMortuum/servus/calendar/validate"
	"github.com/DictumMortuum/servus/config"
	"github.com/DictumMortuum/servus/links"
	"github.com/DictumMortuum/servus/util"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"log"
	"os"
)

func main() {
	mode := os.Getenv("GIN_MODE")
	templates := "templates/*"

	if mode == "release" {
		gin.DisableConsoleColor()
		f, _ := os.Create("/var/log/servus.log")
		gin.DefaultWriter = io.MultiWriter(f)
		templates = "/usr/share/servus/*"
	}

	err := config.Read("/etc/servusrc")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"formatDate":       util.FormatDate,
		"formatShift":      calendar.FormatShift,
		"formatShiftColor": calendar.FormatShiftColor,
	})

	r.LoadHTMLGlob(templates)

	cal := r.Group("/calendar")
	{
		cal.GET("/generate", generate.Handler)
		cal.POST("/parse", parse.Handler)
		cal.POST("/validate", validate.Validate)
		cal.GET("/", parse.Render)
	}

	r.POST("/links", links.Handler)
	r.Run("127.0.0.1:1234")
}
