package search

import (
	// "github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	// "github.com/agnivade/levenshtein"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	store "github.com/velebak/colly-sqlite3-storage/colly/sqlite3"
	"log"
	// "net/url"
	// "regexp"
	// "fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lithammer/fuzzysearch/fuzzy"
	// "github.com/remeh/sizedwaitgroup"
	"sort"
	"strings"
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
}

func initializeScraper(pwd string) (*colly.Collector, *queue.Queue, *store.Storage, error) {
	storage := &store.Storage{
		Filename: pwd + "/results.db",
	}

	collector := colly.NewCollector(
		colly.AllowedDomains(allowed_domains...),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)
	collector.CacheDir = pwd

	err := collector.SetStorage(storage)
	if err != nil {
		return nil, nil, nil, err
	}

	queue, err := queue.New(2, storage)
	if err != nil {
		return nil, nil, nil, err
	}

	return collector, queue, storage, nil
}

func getBoardgameNames(db *sqlx.DB) ([]string, error) {
	var rs []string

	sql := `
		select
			name
		from
			tboardgames
	`

	err := db.Select(&rs, sql)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func boardgameNameToId(db *sqlx.DB, name string) (*models.Boardgame, error) {
	var rs models.Boardgame

	sql := `
		select * from tboardgames where name = ?
	`

	err := db.QueryRowx(sql, name).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func fuzzyFind(col []string) func(string) fuzzy.Ranks {
	return func(s string) fuzzy.Ranks {
		rs := fuzzy.RankFindNormalizedFold(s, col)
		sort.Sort(rs)
		l := len(rs)

		hi := 1
		if l < 1 {
			hi = l
		}

		return rs[0:hi]
	}
}

func Scrape(db *sqlx.DB, batch_id *models.JsonNullInt64) ([]models.Price, error) {
	rs := []models.Price{}

	collector, queue, storage, err := initializeScraper("/tmp")
	if err != nil {
		return nil, err
	}
	defer storage.Close()

	collector.OnHTML(".next a", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.kaissa.eu") {
			return
		}

		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML("article.product", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.kaissa.eu") {
			return
		}

		raw_price := e.ChildText(".price")

		rs = append(rs, models.Price{
			Name:       e.ChildText(".caption"),
			StoreId:    6,
			StoreThumb: e.ChildAttr(".photo a img", "src"),
			Stock:      !childHasClass(e, ".add-to-cart input", "stock-update"),
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".photo a", "href")),
		})
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.thegamerules.com") {
			return
		}

		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML(".product-layout", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.thegamerules.com") {
			return
		}

		raw_price := e.ChildText(".price-new")

		if raw_price == "" {
			raw_price = e.ChildText(".price-normal")
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    4,
			StoreThumb: e.ChildAttr(".product-img div img", "data-src"),
			Stock:      !hasClass(e, "out-of-stock"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		})
	})

	collector.OnHTML(".woocommerce-pagination a.next", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "store.v-games.gr") {
			return
		}

		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML("li.product.type-product", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "store.v-games.gr") {
			return
		}

		raw_price := e.ChildText(".price")

		rs = append(rs, models.Price{
			Name:       e.ChildText(".woocommerce-loop-product__title"),
			StoreId:    5,
			StoreThumb: e.ChildAttr(".woocommerce-LoopProduct-link.woocommerce-loop-product__link img", "data-src"),
			Stock:      childHasClass(e, ".button.product_type_simple", "add_to_cart_button"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".woocommerce-LoopProduct-link", "href"),
		})
	})

	collector.OnHTML("a.next.page-number", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "meeple-planet.com") {
			return
		}

		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML("div.product-small.product-type-simple", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "meeple-planet.com") {
			return
		}

		raw_price := e.ChildText(".amount")

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    7,
			StoreThumb: e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			Stock:      !hasClass(e, "out-of-stock"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		})
	})

	collector.OnHTML("a.--next", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.efantasy.gr") {
			return
		}

		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML("div.product.product-box", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.efantasy.gr") {
			return
		}

		raw_price := e.ChildText(".product-price a strong")

		rs = append(rs, models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    8,
			StoreThumb: e.ChildAttr(".product-image a img", "src"),
			Stock:      !(e.ChildAttr(".product-actions button.cartbutton", "data-stock") == "0"),
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".product-title a", "href")),
		})
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "kaissagames.com") {
			return
		}

		link := e.Attr("href")
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML("li.item.product-item", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "kaissagames.com") {
			return
		}

		raw_price := e.ChildText(".price")

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    9,
			StoreThumb: e.ChildAttr(".product-image-photo", "src"),
			Stock:      !childHasClass(e, "div.stock", "unavailable"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		})
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "meepleonboard.gr") {
			return
		}

		link := e.Attr("href")
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML("div.product-small.product-type-simple", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "meepleonboard.gr") {
			return
		}

		raw_price := e.ChildText(".amount")

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    10,
			StoreThumb: e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			Stock:      !hasClass(e, "out-of-stock"),
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
		})
	})

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

	queue.AddURL("https://www.kaissa.eu/products/epitrapezia-kaissa")
	queue.AddURL("https://www.kaissa.eu/products/epitrapezia-sta-agglika")
	queue.AddURL("https://www.thegamerules.com/epitrapezia-paixnidia?fa132=Board%20Game%20Expansions,Board%20Games")
	queue.AddURL("https://store.v-games.gr/category/board-games")
	queue.AddURL("https://meeple-planet.com/category/epitrapezia-paixnidia")
	queue.AddURL("https://www.efantasy.gr/en/products/%CE%B5%CF%80%CE%B9%CF%84%CF%81%CE%B1%CF%80%CE%AD%CE%B6%CE%B9%CE%B1-%CF%80%CE%B1%CE%B9%CF%87%CE%BD%CE%AF%CE%B4%CE%B9%CE%B1/sc-all")
	queue.AddURL("https://kaissagames.com/b2c_gr/xenoglossa-epitrapezia.html")
	queue.AddURL("https://meepleonboard.gr/product-category/board-games")
	// queue.AddURL("https://www.skroutz.gr/c/259/epitrapezia/shop/7101/No-Label-X.html")
	queue.Run(collector)

	err = storage.Clear()
	if err != nil {
		return nil, err
	}

	for _, item := range rs {
		item.BoardgameId = models.JsonNullInt64{
			Int64: -1,
			Valid: false,
		}

		exists, err := exists(db, item)
		if err != nil {
			return nil, err
		}

		if !exists {
			_, err := create(db, item)
			if err != nil {
				return nil, err
			}
		}
	}

	// 1
	// for i, item := range rs {
	// 	results := fz(item.Name)
	// 	l := len(rs)
	// 	if len(results) > 0 {
	// 		//fmt.Printf("[%5v/%v] %s -----> %v\n", i, l, item.Name, fz(item.Name))

	// 		g, err := boardgameNameToId(db, item.Name)
	// 		if err != nil {
	// 			continue
	// 		}

	// 		fmt.Printf("[%5v/%v] %v to %v\n", i, l, results[0], g.Id)

	// 		item.BoardgameId = g.Id

	// 		if batch_id != nil {
	// 			item.Batch = batch_id.Int64
	// 		}

	// 		item.Hamming = Hamming(item.Name, results[0].Target)
	// 		item.Levenshtein = results[0].Distance

	// 		if batch_id != nil {
	// 			_, err := create(db, item)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 		}
	// 	}
	// }

	// 2
	// swg := sizedwaitgroup.New(8)
	// l := len(rs)

	// for i := 0; i < l; i++ {
	// 	swg.Add()

	// 	go func() {
	// 		defer swg.Done()
	// 		results := fz(rs[i].Name)

	// 		if len(results) > 0 {
	// 			g, err := boardgameNameToId(db, rs[i].Name)
	// 			if err != nil {
	// 				return
	// 			}

	// 			fmt.Printf("[%5v/%v] %v to %v\n", i, l, results[0], g.Id)

	// 			rs[i].BoardgameId = g.Id

	// 			if batch_id != nil {
	// 				rs[i].Batch = batch_id.Int64
	// 			}

	// 			rs[i].Hamming = Hamming(rs[i].Name, results[0].Target)
	// 			rs[i].Levenshtein = results[0].Distance

	// 			if batch_id != nil {
	// 				_, err := create(db, rs[i])
	// 				if err != nil {
	// 					fmt.Println("[%5v/%v] %s", err)
	// 					return
	// 				}
	// 			}
	// 		}
	// 	}()
	// }

	// swg.Wait()

	return rs, nil
}
