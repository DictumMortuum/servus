package main

import (
	"github.com/DictumMortuum/servus/boardgames"
	"github.com/DictumMortuum/servus/calendar"
	"github.com/DictumMortuum/servus/calendar/generate"
	"github.com/DictumMortuum/servus/calendar/parse"
	"github.com/DictumMortuum/servus/calendar/validate"
	"github.com/DictumMortuum/servus/config"
	"github.com/DictumMortuum/servus/gas"
	"github.com/DictumMortuum/servus/gnucash"
	// "github.com/DictumMortuum/servus/generic"
	"github.com/DictumMortuum/servus/links"
	"github.com/DictumMortuum/servus/music"
	"github.com/DictumMortuum/servus/router"
	"github.com/DictumMortuum/servus/util"
	"github.com/DictumMortuum/servus/weight"
	"github.com/DictumMortuum/servus/zerotier"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"html/template"
	"io"
	"log"
	"os"
	"time"
)

func Version(c *gin.Context) {
	util.Success(c, map[string]string{
		"version": "1.2.1",
	})
}

// SetConfig gin Middlware to push some config values
func SetConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("CorsOrigin", "*")
		c.Set("Verbose", true)
		c.Next()
	}
}

// Options common response for rest options
func Options(c *gin.Context) {
	Origin := c.MustGet("CorsOrigin").(string)

	c.Writer.Header().Set("Access-Control-Allow-Origin", Origin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

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
	r.Use(SetConfig())
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, Bearer, range",
		ExposedHeaders:  "x-total-count, Content-Range",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

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
		cal.GET("/sync", parse.SyncToNextcloud)
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
		bg.GET("/prices/notify", boardgames.SendNotifications)
		bg.GET("/duel", func(c *gin.Context) {
			rs, err := boardgames.GetDuel()
			if err != nil {
				util.Error(c, err)
				return
			}
			util.Success(c, &rs)
		})
		bg.GET("/wingspan", func(c *gin.Context) {
			rs, err := boardgames.GetWingspan()
			if err != nil {
				util.Error(c, err)
				return
			}
			util.Success(c, &rs)
		})
	}

	gn := r.Group("/gnucash")
	{
		gn.GET("/expenses/:expense", gnucash.GetExpenseByMonth)
	}

	rt := r.Group("/router")
	{
		rt.GET("", router.Get)
		rt.GET("/latest", router.Latest)
	}

	r.POST("/weight", weight.AddWeight)
	r.POST("/links", links.AddLink)
	r.GET("/version", Version)
	r.Run("127.0.0.1:1234")
}
