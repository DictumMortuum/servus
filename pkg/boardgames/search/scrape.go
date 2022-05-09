package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/jmoiron/sqlx"
	"log"
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
	"epitrapez.io",
	"www.ozon.gr",
	"gamesuniverse.gr",
	"xrysoftero.gr",
	"www.gameexplorers.gr",
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
	rs := []models.Price{}

	collector, queue, err := initializeScraper("/tmp")
	if err != nil {
		return nil, err
	}

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

		var stock int
		raw_price := e.ChildText(".price")

		if childHasClass(e, ".add-to-cart input", "stock-update") {
			stock = 2
		} else {
			stock = 0
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".caption"),
			StoreId:    6,
			StoreThumb: e.ChildAttr(".photo a img", "src"),
			Stock:      stock,
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

		var stock int

		if !hasClass(e, "out-of-stock") {
			stock = 0
		} else {
			stat := e.ChildText(".stat-1")

			if stat == "Preorder" {
				stock = 1
			} else {
				stock = 2
			}
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    4,
			StoreThumb: e.ChildAttr(".product-img div img", "data-src"),
			Stock:      stock,
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

		var stock int

		if hasClass(e, "instock") {
			stock = 0
		} else if hasClass(e, "onbackorder") {
			stock = 1
		} else if hasClass(e, "outofstock") {
			stock = 2
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".woocommerce-loop-product__title"),
			StoreId:    5,
			StoreThumb: e.ChildAttr(".woocommerce-LoopProduct-link.woocommerce-loop-product__link img", "data-src"),
			Stock:      stock,
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

		var stock int

		if hasClass(e, "out-of-stock") {
			stock = 2
		} else {
			if e.ChildText(".badge-inner") != "" {
				stock = 1
			} else {
				stock = 0
			}
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    7,
			StoreThumb: e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		})
	})

	collector.OnHTML(".pagination a", func(e *colly.HTMLElement) {
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

		var stock int

		if e.Attr("data-label") == "" {
			stock = 2
		} else if e.Attr("data-label") == "preorder" {
			stock = 1
		} else {
			stock = 0
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    8,
			StoreThumb: e.ChildAttr(".product-image a img", "src"),
			Stock:      stock,
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

		var stock int

		if e.ChildText(".release-date") != "" {
			stock = 1
		} else {
			if !childHasClass(e, "div.stock", "unavailable") {
				stock = 0
			} else {
				stock = 2
			}
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    9,
			StoreThumb: e.ChildAttr(".product-image-photo", "src"),
			Stock:      stock,
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

		var stock int

		if hasClass(e, "instock") {
			stock = 0
		} else if hasClass(e, "onbackorder") {
			stock = 1
		} else if hasClass(e, "out-of-stock") {
			stock = 2
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    10,
			StoreThumb: e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
		})
	})

	collector.OnHTML("._3DNsL", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "mystery-bay") {
			return
		}

		raw_price := e.ChildText("._2-l9W")
		raw_url := e.ChildAttr("._1FMIK", "style")
		urls := getURL(raw_url)

		url := ""
		if len(urls) > 0 {
			url = urls[0]
		}

		var stock int

		if e.ChildText("span[data-hook=product-item-ribbon-new]") == "PRE-ORDER" {
			stock = 1
		} else {
			if e.ChildAttr("button[data-hook=product-item-add-to-cart-button]", "aria-disabled") == "true" {
				stock = 2
			} else {
				stock = 0
			}
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText("h3"),
			StoreId:    3,
			StoreThumb: url,
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr("a[data-hook=product-item-container]", "href"),
		})
	})

	collector.OnHTML(".woocommerce-pagination a.next", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "https://epitrapez.io") {
			return
		}

		link := e.Attr("href")
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML("li.product.type-product", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "https://epitrapez.io") {
			return
		}

		raw_price := e.ChildText(".price ins .amount")

		if raw_price == "" {
			raw_price = e.ChildText(".price .amount")
		}

		var stock int

		if e.ChildText("a.add_to_cart_button") != "" {
			stock = 0
		} else {
			stock = 2
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".woocommerce-loop-product__title"),
			StoreId:    15,
			StoreThumb: e.ChildAttr(".epz-product-thumbnail img", "data-src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".woocommerce-LoopProduct-link", "href"),
		})
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "https://www.ozon.gr") {
			return
		}

		link := e.Attr("href")
		if link != "javascript:;" {
			log.Println("Visiting: " + link)
			queue.AddURL(link)
		}
	})

	collector.OnHTML(".products-list div.col-xs-3", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "https://www.ozon.gr") {
			return
		}

		raw_price := e.ChildText(".special-price")

		if raw_price == "" {
			raw_price = e.ChildText(".price")
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".title"),
			StoreId:    17,
			StoreThumb: e.ChildAttr(".image-wrapper img", "src"),
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".product-box", "href"),
		})
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "gamesuniverse.gr") {
			return
		}

		link := e.Attr("href")
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML("article.product-miniature", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "gamesuniverse.gr") {
			return
		}

		raw_price := e.ChildText(".product-price")

		rs = append(rs, models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    20,
			StoreThumb: e.ChildAttr(".thumbnail img", "data-src"),
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".product-thumbnail", "href"),
		})
	})

	collector.OnHTML("a.pagenav", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.politeianet.gr") {
			return
		}

		link := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML(".browse-page-block", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.politeianet.gr") {
			return
		}

		raw_price := e.ChildText(".productPrice")
		if raw_price == "" {
			return
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".browse-product-title"),
			StoreId:    12,
			StoreThumb: e.ChildAttr(".browseProductImage", "src"),
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".browse-product-title", "href"),
		})
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "xrysoftero.gr") {
			return
		}

		link := e.Attr("href")
		log.Println("Visiting: " + link)
		queue.AddURL(link)
	})

	collector.OnHTML(".thumbnail-container", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "xrysoftero.gr") {
			return
		}

		url := e.ChildAttr(".cover-image", "src")
		if url == "" {
			return
		}

		raw_price := e.ChildText(".price")
		rs = append(rs, models.Price{
			Name:       e.ChildText(".product-title"),
			StoreId:    21,
			StoreThumb: url,
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr("a.relative", "href"),
		})
	})

	collector.OnHTML(".product-pagination > a", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.gameexplorers.gr") {
			return
		}

		if e.Attr("title") == "επόμενη σελίδα" {
			link := e.Attr("href")
			log.Println("Visiting: " + link)
			queue.AddURL(link)
		}
	})

	collector.OnHTML(".single-product-item", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.gameexplorers.gr") {
			return
		}

		raw_price := e.ChildText(".regular-price")
		rs = append(rs, models.Price{
			Name:       e.ChildText("h2:nth-child(1)"),
			StoreId:    22,
			StoreThumb: e.ChildAttr("a:nth-child(1) > img:nth-child(1)", "src"),
			Stock:      0,
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr("a:nth-child(1)", "href"),
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
	// queue.AddURL("https://www.skroutz.gr/c/259/epitrapezia/shop/7101/No-Label-X.html")

	queue.AddURL("https://www.kaissa.eu/products/epitrapezia-kaissa")
	queue.AddURL("https://www.kaissa.eu/products/epitrapezia-sta-agglika")
	queue.AddURL("https://www.thegamerules.com/epitrapezia-paixnidia?fa132=Board%20Game%20Expansions,Board%20Games")
	queue.AddURL("https://store.v-games.gr/category/board-games")
	queue.AddURL("https://meeple-planet.com/category/epitrapezia-paixnidia")
	queue.AddURL("https://meeple-planet.com/category/pre-orders")
	queue.AddURL("https://www.efantasy.gr/en/products/%CE%B5%CF%80%CE%B9%CF%84%CF%81%CE%B1%CF%80%CE%AD%CE%B6%CE%B9%CE%B1-%CF%80%CE%B1%CE%B9%CF%87%CE%BD%CE%AF%CE%B4%CE%B9%CE%B1/sc-all")
	queue.AddURL("https://kaissagames.com/b2c_gr/xenoglossa-epitrapezia.html")
	queue.AddURL("https://meepleonboard.gr/product-category/board-games")
	queue.AddURL("https://epitrapez.io/product-category/epitrapezia/?Stock=allstock")
	queue.AddURL("https://www.mystery-bay.com/epitrapezia-paixnidia?page=36")
	queue.AddURL("https://www.ozon.gr/pazl-kai-paixnidia/epitrapezia-paixnidia")
	queue.AddURL("https://gamesuniverse.gr/el/10-epitrapezia")
	queue.AddURL("https://www.politeianet.gr/index.php?option=com_virtuemart&category_id=948&page=shop.browse&subCatFilter=-1&langFilter=-1&pubdateFilter=-1&availabilityFilter=-1&discountFilter=-1&priceFilter=-1&pageFilter=-1&Itemid=721&limit=20&limitstart=0")
	queue.AddURL("https://xrysoftero.gr/362-epitrapezia-paixnidia?resultsPerPage=48&q=%CE%9C%CE%AC%CF%81%CE%BA%CE%B1%5C-%CE%95%CE%BA%CE%B4%CF%8C%CF%84%CE%B7%CF%82-%CE%9A%CE%AC%CE%B9%CF%83%CF%83%CE%B1")
	queue.AddURL("https://www.gameexplorers.gr/kartes-epitrapezia/epitrapezia-paixnidia/items-grid-date-desc-1-60.html")
	queue.Run(collector)

	store_ids := []int64{}

	for _, item := range rs {
		store_ids = append(store_ids, item.StoreId)
	}

	store_ids = unique(store_ids)

	for _, store_id := range store_ids {
		err := updateBatch(db, store_id)
		if err != nil {
			return nil, err
		}
	}

	for _, item := range rs {
		err = upsertPrice(db, item)
		if err != nil {
			return nil, err
		}
	}

	return rs, nil
}
