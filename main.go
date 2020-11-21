package main

import (
	"github.com/DictumMortuum/servus/calendar"
	"github.com/DictumMortuum/servus/calendar/generate"
	"github.com/DictumMortuum/servus/calendar/parse"
	"github.com/DictumMortuum/servus/calendar/validate"
	"github.com/DictumMortuum/servus/config"
	"github.com/DictumMortuum/servus/gas"
	"github.com/DictumMortuum/servus/links"
	"github.com/DictumMortuum/servus/util"
	"github.com/DictumMortuum/servus/zerotier"
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
		templates = "/opt/domus/servus/templates/*"
	}

	err := config.Read("/etc/servusrc")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"formatDate":       util.FormatDate,
		"formatDay":        util.FormatDay,
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

	zt := r.Group("/zerotier")
	{
		zt.GET("/member/:member", zerotier.PostNode)
	}

	gs := r.Group("/fuel")
	{
		gs.GET("/", gas.Render)
		gs.POST("/addstats", gas.AddFuelStats)
		gs.POST("/add", gas.AddFuel)
	}

	r.POST("/links", links.AddLink)
	r.Run("0.0.0.0:1234")
}
