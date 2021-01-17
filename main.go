package main

import (
	"github.com/DictumMortuum/servus/calendar"
	"github.com/DictumMortuum/servus/calendar/generate"
	"github.com/DictumMortuum/servus/calendar/parse"
	"github.com/DictumMortuum/servus/calendar/validate"
	"github.com/DictumMortuum/servus/config"
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/gas"
	"github.com/DictumMortuum/servus/generic"
	"github.com/DictumMortuum/servus/links"
	"github.com/DictumMortuum/servus/util"
	"github.com/DictumMortuum/servus/zerotier"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"html/template"
	"io"
	"log"
	"os"
	"time"
)

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
		zt.GET("/member/:member", zerotier.GetNode)
		zt.POST("/member", zerotier.PostNode)
	}

	gs := r.Group("/fuel")
	{
		gs.GET("/", gas.Render)
		gs.POST("/addstats", gas.AddFuelStats)
		gs.POST("/add", gas.AddFuel)
	}

	gen := r.Group("/generic")
	{
		gen.GET("", db.GetShiftList)
		gen.GET("/getList", db.GetShiftList)
		gen.GET("/getOne/:id", generic.GetOne(generic.GetOneSample))
		gen.POST("", generic.Create(generic.CreateSample))
		gen.DELETE("/:id", generic.Delete(generic.DeleteSample))
		gen.OPTIONS("", Options)
	}

	r.POST("/links", links.AddLink)
	r.Run("0.0.0.0:1234")
}
