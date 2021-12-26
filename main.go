package main

import (
	"embed"
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
	"github.com/DictumMortuum/servus/pkg/tasks"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/DictumMortuum/servus/pkg/weight"
	"github.com/DictumMortuum/servus/pkg/zerotier"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"html/template"
	"io"
	"log"
	"net/http"
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

//go:embed assets
var staticFS embed.FS

func main() {
	mode := os.Getenv("GIN_MODE")
	path_templates := "templates/*"
	path_cfg := "servusrc"

	if mode == "release" {
		gin.DisableConsoleColor()
		f, _ := os.Create("/var/log/servus.log")
		gin.DefaultWriter = io.MultiWriter(f)
		path_templates = "/usr/share/webapps/servus/*"
		path_cfg = "/etc/servusrc"
	}

	err := config.Read(path_cfg)
	if err != nil {
		return
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

	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})

	r.GET("/version", func(c *gin.Context) {
		c.FileFromFS("/assets/version.json", http.FS(staticFS))
	})

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
		ms.GET("/current", music.Current)
		ms.GET("/radio", music.Radio)
	}

	bg := r.Group("/boardgames")
	{
		bg.GET("/prices/notify", prices.SendMessages)
		bg.GET("/prices/msg", prices.CreateMessages)
		bg.GET("/prices/new", prices.GetUpdates)
		bg.POST("/search/:hash", CacheCheck(apiCache), boardgames.AtlasSearch)
		bg.POST("/get/:hash", CacheCheck(apiCache), boardgames.BggGet)
		bg.POST("/bggsearch/:hash", CacheCheck(apiCache), boardgames.BggSearch)
	}

	ts := r.Group("/tasks")
	{
		ts.GET("/:list", tasks.GetTasks)
		ts.GET("/:list/sync", tasks.SyncTasks)
	}

	rest := r.Group("/rest/v1")
	{
		rest.GET("/boardgame/:id", generic.F(boardgames.GetBoardgame))
		rest.GET("/boardgame", generic.F(boardgames.GetListBoardgame))
		rest.POST("/boardgame", generic.F(boardgames.CreateBoardgame))
		rest.PUT("/boardgame/:id", generic.F(boardgames.UpdateBoardgame))
		rest.DELETE("/boardgame/:id", generic.F(boardgames.DeleteBoardgame))

		rest.GET("/store/:id", generic.F(boardgames.GetStore))
		rest.GET("/store", generic.F(boardgames.GetListStore))
		rest.POST("/store", generic.F(boardgames.CreateStore))
		rest.PUT("/store/:id", generic.F(boardgames.UpdateStore))
		rest.DELETE("/store/:id", generic.F(boardgames.DeleteStore))

		rest.GET("/scrape/:id", generic.F(scraper.GetData))
		rest.GET("/scrape", generic.F(scraper.GetListData))
		rest.POST("/scrape", generic.F(scraper.CreateData))
		rest.PUT("/scrape/:id", generic.F(scraper.UpdateData))
		rest.DELETE("/scrape/:id", generic.F(scraper.DeleteData))

		rest.GET("/player/:id", generic.F(boardgames.GetPlayer))
		rest.GET("/player", generic.F(boardgames.GetListPlayer))
		rest.POST("/player", generic.F(boardgames.CreatePlayer))
		rest.PUT("/player/:id", generic.F(boardgames.UpdatePlayer))
		rest.DELETE("/player/:id", generic.F(boardgames.DeletePlayer))

		rest.GET("/play/:id", generic.F(boardgames.GetPlay))
		rest.GET("/play", generic.F(boardgames.GetListPlay))
		rest.POST("/play", generic.F(boardgames.CreatePlay))
		rest.PUT("/play/:id", generic.F(boardgames.UpdatePlay))
		rest.DELETE("/play/:id", generic.F(boardgames.DeletePlay))

		rest.GET("/stats/:id", generic.F(boardgames.GetStats))
		rest.GET("/stats", generic.F(boardgames.GetListStats))
		rest.POST("/stats", generic.F(boardgames.CreateStats))
		rest.PUT("/stats/:id", generic.F(boardgames.UpdateStats))
		rest.DELETE("/stats/:id", generic.F(boardgames.DeleteStats))
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
	r.GET("/expenses", gnucash.GetTopExpenses)
	r.GET("/cache", CacheSave(apiCache))
	r.Run("127.0.0.1:1234")
}
