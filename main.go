package main

import (
	"github.com/DictumMortuum/servus/boardgames"
	"github.com/DictumMortuum/servus/calendar"
	"github.com/DictumMortuum/servus/calendar/generate"
	"github.com/DictumMortuum/servus/calendar/parse"
	"github.com/DictumMortuum/servus/calendar/validate"
	"github.com/DictumMortuum/servus/config"
	"github.com/DictumMortuum/servus/gas"
	"github.com/DictumMortuum/servus/links"
	"github.com/DictumMortuum/servus/music"
	"github.com/DictumMortuum/servus/router"
	"github.com/DictumMortuum/servus/util"
	"github.com/DictumMortuum/servus/weight"
	"github.com/DictumMortuum/servus/zerotier"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"log"
	"os"
)

func main() {
	mode := os.Getenv("GIN_MODE")
	path_templates := "templates/*"
	path_cfg := "servusrc"

	if mode == "release" {
		gin.DisableConsoleColor()
		f, _ := os.Create("/var/log/servus.log")
		gin.DefaultWriter = io.MultiWriter(f)
		path_templates = "/opt/domus/servus/templates/*"
		path_cfg = "/etc/servusrc"
	}

	err := config.Read(path_cfg)
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

	r.LoadHTMLGlob(path_templates)

	cal := r.Group("/calendar")
	{
		cal.GET("/generate", generate.Handler)
		cal.POST("/parse", parse.Handler)
		cal.POST("/validate", validate.Validate)
		cal.GET("/", parse.Render)
	}

	zt := r.Group("/zerotier")
	{
		zt.GET("/member/:member", zerotier.GetNode)
		zt.POST("/member", zerotier.PostNode)
	}

	gs := r.Group("/fuel")
	{
		gs.GET("/", gas.Render)
		gs.POST("/addstats", gas.AddFuelStats)
		gs.POST("/add", gas.AddFuel)
	}

	ms := r.Group("/music")
	{
		ms.GET("/playlist/:playlist", music.Playlist)
		ms.GET("/stop", music.Stop)
	}

	bg := r.Group("/boardgames")
	{
		bg.GET("/prices", boardgames.GetPrices)
	}

	r.GET("/router", router.Get)
	r.POST("/weight", weight.AddWeight)
	r.POST("/links", links.AddLink)
	r.Run("127.0.0.1:1234")
}
