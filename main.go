package main

import (
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/calendar"
	"github.com/DictumMortuum/servus/pkg/calendar/generate"
	"github.com/DictumMortuum/servus/pkg/calendar/parse"
	"github.com/DictumMortuum/servus/pkg/calendar/validate"
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/gas"
	"github.com/DictumMortuum/servus/pkg/generic"
	"github.com/DictumMortuum/servus/pkg/gnucash"
	"github.com/DictumMortuum/servus/pkg/links"
	"github.com/DictumMortuum/servus/pkg/music"
	"github.com/DictumMortuum/servus/pkg/prices"
	"github.com/DictumMortuum/servus/pkg/router"
	"github.com/DictumMortuum/servus/pkg/scraper"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/DictumMortuum/servus/pkg/weight"
	"github.com/DictumMortuum/servus/pkg/zerotier"
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
		"version": "3.0.1",
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

	apiCache, err := CacheInit()
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
		bg.GET("/prices/notify", prices.SendMessages)
		bg.GET("/prices/msg", prices.CreateMessages)
		bg.GET("/prices/new", prices.GetUpdates)
		bg.POST("/search/:hash", CacheCheck(apiCache), boardgames.AtlasSearch)
		bg.POST("/get/:hash", CacheCheck(apiCache), boardgames.BggGet)
		bg.POST("/bggsearch/:hash", CacheCheck(apiCache), boardgames.BggSearch)
		bg.GET("/scores", boardgames.GetScores)
	}

	rest := r.Group("/rest/v1")
	{
		games := boardgames.Boardgame{}
		rest.GET("/boardgame/:id", generic.GET(games))
		rest.GET("/boardgame", generic.GETLIST(games))
		rest.POST("/boardgame", generic.POST(games))
		rest.PUT("/boardgame/:id", generic.PUT(games))
		rest.DELETE("/boardgame/:id", generic.DELETE(games))

		store := boardgames.Store{}
		rest.GET("/store/:id", generic.GET(store))
		rest.GET("/store", generic.GETLIST(store))
		rest.POST("/store", generic.POST(store))
		rest.PUT("/store/:id", generic.PUT(store))
		rest.DELETE("/store/:id", generic.DELETE(store))

		data := scraper.Data{}
		rest.GET("/scrape/:id", generic.GET(data))
		rest.GET("/scrape", generic.GETLIST(data))
		rest.POST("/scrape", generic.POST(data))
		rest.PUT("/scrape/:id", generic.PUT(data))
		rest.DELETE("/scrape/:id", generic.DELETE(data))

		player := boardgames.Player{}
		rest.GET("/player/:id", generic.GET(player))
		rest.GET("/player", generic.GETLIST(player))
		rest.POST("/player", generic.POST(player))
		rest.PUT("/player/:id", generic.PUT(player))
		rest.DELETE("/player/:id", generic.DELETE(player))

		play := boardgames.Play{}
		rest.GET("/play/:id", generic.GET(play))
		rest.GET("/play", generic.GETLIST(play))
		rest.POST("/play", generic.POST(play))
		rest.PUT("/play/:id", generic.PUT(play))
		rest.DELETE("/play/:id", generic.DELETE(play))

		stats := boardgames.Stats{}
		rest.GET("/stats/:id", generic.GET(stats))
		rest.GET("/stats", generic.GETLIST(stats))
		rest.POST("/stats", generic.POST(stats))
		rest.PUT("/stats/:id", generic.PUT(stats))
		rest.DELETE("/stats/:id", generic.DELETE(stats))
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

	sc := r.Group("/scrape")
	{
		sc.POST("/data", scraper.CreateDataMapping)

		vgames := scraper.VgamesScraper{"V Games"}
		sc.GET("/vgames", scraper.Scrape(vgames))
		sc.GET("/vgames/prices", scraper.ScrapePrices(vgames))

		gamerules := scraper.GameRulesScraper{"The Game Rules"}
		sc.GET("/gamerules", scraper.Scrape(gamerules))
		sc.GET("/gamerules/prices", scraper.ScrapePrices(gamerules))

		mysterybay := scraper.MysteryBayScraper{"Mystery Bay"}
		sc.GET("/mystery", scraper.Scrape(mysterybay))
		sc.GET("/mystery/prices", scraper.ScrapePrices(mysterybay))

		kaissa := scraper.KaissaScraper{"Kaissa Amarousiou"}
		sc.GET("/kaissa", scraper.Scrape(kaissa))
		sc.GET("/kaissa/prices", scraper.ScrapePrices(kaissa))

		fantasygate := scraper.FantasyGateScraper{"Fantasy Gate"}
		sc.GET("/fantasygate", scraper.Scrape(fantasygate))
		sc.GET("/fantasygate/prices", scraper.ScrapePrices(fantasygate))
	}

	r.POST("/weight", weight.AddWeight)
	r.POST("/links", links.AddLink)
	r.GET("/version", Version)
	r.GET("/cache", CacheSave(apiCache))
	r.Run("127.0.0.1:1234")
}
