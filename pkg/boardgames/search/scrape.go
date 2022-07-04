package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/rabbitmq"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/jmoiron/sqlx"
	"log"
)

var allowed_domains = []string{
	"www.thegamerules.com",
	"www.mystery-bay.com",
	"store.v-games.gr",
	"www.kaissa.eu",
	"meeple-planet.com",
	"kaissagames.com",
	"www.efantasy.gr",
	"meepleonboard.gr",
	"www.gameshero.gr",
	"www.politeianet.gr",
	"www.skroutz.gr",
	"epitrapez.io",
	"www.ozon.gr",
	"gamesuniverse.gr",
	"xrysoftero.gr",
	"www.gameexplorers.gr",
	"crystallotus.eu",
	"avalongames.gr",
	"rollnplay.gr",
}

func initializeScraper(pwd string) (*colly.Collector, *queue.Queue, error) {
	collector := colly.NewCollector(
		colly.AllowedDomains(allowed_domains...),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)
	collector.CacheDir = pwd

	queue, err := queue.New(2, &queue.InMemoryQueueStorage{MaxSize: 100000})
	if err != nil {
		return nil, nil, err
	}

	return collector, queue, nil
}

func Scrape(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	conn, ch, q, err := rabbitmq.SetupQueue("prices")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer ch.Close()

	log.Println(q)

	store_ids := []int64{6, 5, 7, 4, 9, 8, 10, 15, 17, 3, 20, 12, 21, 22, 24, 25, 26}
	for _, store_id := range store_ids {
		err := updateBatch(db, store_id)
		if err != nil {
			return nil, err
		}
	}

	collector, queue, err := initializeScraper("/tmp")
	if err != nil {
		return nil, err
	}

	// // No Label X
	// collector.OnHTML("a.next", func(e *colly.HTMLElement) {
	// 	if !strings.Contains(e.Request.URL.String(), "www.skroutz.gr") {
	// 		return
	// 	}

	// 	link := e.Request.AbsoluteURL(e.Attr("href"))
	// 	log.Println("Visiting: " + link)
	// 	queue.AddURL(link)
	// })

	// collector.OnHTML("li.cf.card", func(e *colly.HTMLElement) {
	// 	if !strings.Contains(e.Request.URL.String(), "www.skroutz.gr") {
	// 		return
	// 	}

	// 	raw_price := e.ChildText(".price")

	// 	rs = append(rs, models.Price{
	// 		BoardgameId: -1,
	// 		Name:        e.ChildText(".js-sku-link"),
	// 		StoreId:     10,
	// 		StoreThumb:  e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
	// 		Stock:       !hasClass(e, "out-of-stock"),
	// 		Price:       getPrice(raw_price),
	// 		Url:         e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
	// 	})
	// })
	// queue.AddURL("https://www.skroutz.gr/c/259/epitrapezia/shop/7101/No-Label-X.html")

	// queue.AddURL("https://www.kaissa.eu/products/epitrapezia-kaissa")
	// queue.AddURL("https://www.kaissa.eu/products/epitrapezia-sta-agglika")
	// queue.AddURL("https://www.thegamerules.com/epitrapezia-paixnidia?fa132=Board%20Game%20Expansions,Board%20Games")
	// queue.AddURL("https://store.v-games.gr/category/board-games")
	// queue.AddURL("https://meeple-planet.com/category/epitrapezia-paixnidia")
	// queue.AddURL("https://meeple-planet.com/category/pre-orders")
	// queue.AddURL("https://www.efantasy.gr/en/products/%CE%B5%CF%80%CE%B9%CF%84%CF%81%CE%B1%CF%80%CE%AD%CE%B6%CE%B9%CE%B1-%CF%80%CE%B1%CE%B9%CF%87%CE%BD%CE%AF%CE%B4%CE%B9%CE%B1/sc-all")
	// queue.AddURL("https://kaissagames.com/b2c_gr/xenoglossa-epitrapezia.html")
	// queue.AddURL("https://meepleonboard.gr/product-category/board-games")
	// queue.AddURL("https://epitrapez.io/product-category/epitrapezia/?Stock=allstock")
	// queue.AddURL("https://www.mystery-bay.com/epitrapezia-paixnidia?page=36")
	// queue.AddURL("https://www.ozon.gr/pazl-kai-paixnidia/epitrapezia-paixnidia")
	// queue.AddURL("https://gamesuniverse.gr/el/10-epitrapezia")
	// queue.AddURL("https://www.politeianet.gr/index.php?option=com_virtuemart&category_id=948&page=shop.browse&subCatFilter=-1&langFilter=-1&pubdateFilter=-1&availabilityFilter=-1&discountFilter=-1&priceFilter=-1&pageFilter=-1&Itemid=721&limit=20&limitstart=0")
	// queue.AddURL("https://xrysoftero.gr/362-epitrapezia-paixnidia?resultsPerPage=48&q=%CE%9C%CE%AC%CF%81%CE%BA%CE%B1%5C-%CE%95%CE%BA%CE%B4%CF%8C%CF%84%CE%B7%CF%82-%CE%9A%CE%AC%CE%B9%CF%83%CF%83%CE%B1")
	// queue.AddURL("https://www.gameexplorers.gr/kartes-epitrapezia/epitrapezia-paixnidia/items-grid-date-desc-1-60.html")
	// queue.AddURL("https://crystallotus.eu/collections/board-games")
	// queue.AddURL("https://avalongames.gr/index.php?route=product/category&path=59&limit=100")
	// queue.AddURL("https://rollnplay.gr/?product_cat=all-categories")
	queue.Run(collector)

	return nil, nil
}
