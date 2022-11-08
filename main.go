package main

import (
	"embed"
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/DictumMortuum/servus/pkg/boardgames/images"
	"github.com/DictumMortuum/servus/pkg/boardgames/mapping"
	"github.com/DictumMortuum/servus/pkg/boardgames/mathtrade"
	"github.com/DictumMortuum/servus/pkg/boardgames/search"
	"github.com/DictumMortuum/servus/pkg/books"
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/food"
	"github.com/DictumMortuum/servus/pkg/gas"
	"github.com/DictumMortuum/servus/pkg/generic"
	"github.com/DictumMortuum/servus/pkg/gnucash"
	"github.com/DictumMortuum/servus/pkg/links"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/music"
	"github.com/DictumMortuum/servus/pkg/router"
	"github.com/DictumMortuum/servus/pkg/tasks"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/DictumMortuum/servus/pkg/weight"
	"github.com/DictumMortuum/servus/pkg/zerotier"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"net/http"
)

//go:embed assets
var staticFS embed.FS

func main() {
	r, err := generic.SetupMainRouter("")
	if err != nil {
		log.Fatal(err)
	}

	r.SetFuncMap(template.FuncMap{
		"formatDate":       util.FormatDate,
		"formatDay":        util.FormatDay,
		"formatShift":      util.FormatShift,
		"formatShiftColor": util.FormatShiftColor,
	})

	r.LoadHTMLGlob(config.App.PathTemplates)

	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})

	r.GET("/version", func(c *gin.Context) {
		c.FileFromFS("/assets/version.json", http.FS(staticFS))
	})

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
		ms.GET("/toggle", music.Toggle)
		ms.GET("/next", music.Next)
		ms.GET("/previous", music.Previous)
		ms.GET("/current", music.Current)
		ms.GET("/radio", music.Radio)
	}

	bg := r.Group("/boardgames")
	{
		// bg.POST("/search/:hash", CacheCheck(apiCache), atlas.AtlasSearch)
		// bg.POST("/get/:hash", CacheCheck(apiCache), bgg.BggGet)
		// bg.POST("/bggsearch/:hash", CacheCheck(apiCache), bgg.BggSearch)
		bg.GET("/top", bgg.GetTopBoardgames)
		bg.GET("/top/art", bgg.FetchTopArt)
	}

	ts := r.Group("/tasks")
	{
		ts.GET("/:list", tasks.GetTasks)
		ts.GET("/:list/sync", tasks.SyncTasks)
	}

	rest := r.Group("/rest/v1")
	{
		rest.GET("/scrape/database", generic.F(search.UpdateMappings))
		rest.GET("/scrape", generic.C([]func(*sqlx.DB, *models.QueryBuilder) (interface{}, error){
			search.ScrapeFantasyShop,
			search.ScrapeAvalon,
			search.ScrapeBoardsOfMadness,
			search.ScrapeCrystalLotus,
			search.ScrapeEfantasy,
			search.ScrapeEpitrapezio,
			search.ScrapeFantasyGate,
			search.ScrapeGameExplorers,
			search.ScrapeGameRules,
			search.ScrapeGamesCom,
			search.ScrapeGamesUniverse,
			search.ScrapeHobbyTheory,
			search.ScrapeKaissaEu,
			search.ScrapeKaissaGames,
			search.ScrapeMeepleOnBoard,
			search.ScrapeMeeplePlanet,
			search.ScrapeMysteryBay,
			search.ScrapeOzon,
			search.ScrapePoliteia,
			search.ScrapeRollnplay,
			search.ScrapeVgames,
			search.ScrapeXrysoFtero,
			search.ScrapeGenx,
			search.ScrapeGreekGuild,
		}))

		rest.GET("/mapping2/all", generic.F(mapping.MapAll))
		rest.GET("/mapping2/bgg", generic.F(mapping.MapAllBgg))
		rest.GET("/mapping2/static", generic.F(mapping.MapAllStatic))
		rest.GET("/mapping2/search", generic.F(mapping.SearchMaps))

		rest.GET("/image/:id", generic.F(images.Boardgame))

		rest.GET("/trueskill", generic.F(boardgames.GetTrueskillLists))
		rest.GET("/trueskill/overall", generic.F(boardgames.GetTrueskillOverall))

		rest.GET("/boardgame/:id", generic.F(boardgames.GetBoardgame))
		rest.GET("/boardgame", generic.F(boardgames.GetListBoardgame))
		rest.POST("/boardgame", generic.F(boardgames.CreateBoardgame))
		rest.PUT("/boardgame/:id", generic.F(boardgames.UpdateBoardgame))
		rest.PUT("/boardgame/:id/refetch", generic.F(boardgames.RefetchBoardgame))
		rest.DELETE("/boardgame/:id", generic.F(boardgames.DeleteBoardgame))

		rest.GET("/store/:id", generic.F(boardgames.GetStore))
		rest.GET("/store", generic.F(boardgames.GetListStore))
		rest.POST("/store", generic.F(boardgames.CreateStore))
		rest.PUT("/store/:id", generic.F(boardgames.UpdateStore))
		rest.DELETE("/store/:id", generic.F(boardgames.DeleteStore))

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

		rest.GET("/mapping/:id", generic.F(boardgames.GetMapping))
		rest.GET("/mapping", generic.F(boardgames.GetListMapping))
		rest.POST("/mapping", generic.F(boardgames.CreateMapping))
		rest.PUT("/mapping/:id", generic.F(boardgames.UpdateMapping))
		rest.DELETE("/mapping/:id", generic.F(boardgames.DeleteMapping))

		rest.GET("/book/:id", generic.F(books.GetBook))
		rest.GET("/book", generic.F(books.GetListBook))
		rest.POST("/book", generic.F(books.CreateBook))
		rest.PUT("/book/:id", generic.F(books.UpdateBook))
		rest.DELETE("/book/:id", generic.F(books.DeleteBook))

		rest.GET("/history/:id", generic.F(boardgames.GetHistoricPrice))
		rest.GET("/history", generic.F(boardgames.GetListHistoricPrice))

		rest.GET("/price/:id", generic.F(boardgames.GetPrice))
		rest.GET("/price", generic.F(boardgames.GetListPrice))
		rest.GET("/price/:id/map/static", generic.F(mapping.MapStatic))
		rest.GET("/price/:id/map/atlas", generic.F(mapping.MapAtlas))
		rest.GET("/price/:id/map/bgg", generic.F(mapping.MapBGG))
		rest.GET("/price/:id/unmap", generic.F(boardgames.UnmapPrice))
		rest.POST("/price", generic.F(boardgames.CreatePrice))
		rest.PUT("/price/:id", generic.F(boardgames.UpdatePrice))
		rest.DELETE("/price/:id", generic.F(boardgames.DeletePrice))

		rest.POST("/search/bgg", generic.F(mapping.SearchBGGTerm))
		rest.POST("/search/atlas", generic.F(mapping.SearchAtlasTerm))

		rest.GET("/expense", generic.S("gnucash", gnucash.GetListExpenses))
		rest.GET("/mathtrade/:id", mathtrade.Analyze)
	}

	gn := r.Group("/gnucash")
	{
		gn.GET("/expenses/:expense", gnucash.GetExpenseByMonth)
	}

	rt := r.Group("/router")
	{
		rt.POST("/reset", router.Reset)
	}

	r.POST("/weight", weight.AddWeight)
	r.POST("/food", food.Scrape)
	r.POST("/links", links.AddLink)
	r.GET("/expenses", gnucash.GetTopExpenses)
	r.Run("127.0.0.1:1234")
}
