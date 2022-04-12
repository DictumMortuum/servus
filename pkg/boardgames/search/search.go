package search

import (
	"database/sql"
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/agnivade/levenshtein"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/jmoiron/sqlx"
	"github.com/velebak/colly-sqlite3-storage/colly/sqlite3"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

func findPrice(db *sqlx.DB, payload models.Price) (*models.JsonNullInt64, error) {
	var id models.JsonNullInt64

	q := `
		select
			id
		from
			tboardgameprices
		where
			store_id = :store_id and
			name = :name
	`

	stmt, err := db.PrepareNamed(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Get(&id, payload)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func exists(db *sqlx.DB, payload models.Price) (bool, error) {
	q := `select 1 from tboardgameprices where store_id = :store_id and name = :name`

	rows, err := db.NamedQuery(q, payload)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func updateBatch(db *sqlx.DB, store_id int64) error {
	q := `
		update
			tboardgameprices
		set
			batch = 0
		where
			store_id = :store_id
	`

	_, err := db.NamedExec(q, map[string]interface{}{
		"store_id": store_id,
	})
	if err != nil {
		return err
	}

	return nil
}

func update(db *sqlx.DB, payload models.Price) (bool, error) {
	q := `
		update
			tboardgameprices
		set
			store_thumb = :store_thumb,
			price = :price,
			stock = :stock,
			url = :url,
			batch = 1,
			cr_date = NOW()
		where
			id = :id
	`

	rs, err := db.NamedExec(q, payload)
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func create(db *sqlx.DB, payload models.Price) (bool, error) {
	q := `
		insert into tboardgameprices (
			boardgame_id,
			name,
			store_id,
			store_thumb,
			price,
			stock,
			url,
			levenshtein,
			hamming,
			batch
		) values (
			:boardgame_id,
			:name,
			:store_id,
			:store_thumb,
			:price,
			:stock,
			:url,
			:levenshtein,
			:hamming,
			1
		)`

	rs, err := db.NamedExec(q, payload)
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func Boardgame(boardgame models.Boardgame, batch_id *models.JsonNullInt64) ([]models.Price, error) {
	rs := []models.Price{}

	enc := url.QueryEscape(boardgame.Name)
	cmd := strings.ToLower(boardgame.Name)

	allowed_domains := []string{
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
	}

	storage := &sqlite3.Storage{
		Filename: "/tmp/results.db",
	}
	defer storage.Close()

	s := colly.NewCollector(
		colly.AllowedDomains(allowed_domains...),
	)
	s.CacheDir = "/tmp/cache"

	err := s.SetStorage(storage)
	if err != nil {
		panic(err)
	}

	q, err := queue.New(len(allowed_domains), storage)
	if err != nil {
		return nil, err
	}

	// The Game Rules
	s.OnHTML(".product-layout", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "thegamerules") {
			return
		}

		raw_price := e.ChildText(".price-new")

		if raw_price == "" {
			raw_price = e.ChildText(".price-normal")
		}

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".name"),
			StoreId:    4,
			StoreThumb: e.ChildAttr(".product-img div img", "data-src"),
			Stock:      !hasClass(e, "out-of-stock"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		})
	})
	q.AddURL("https://www.thegamerules.com/index.php?route=product/search&search=" + enc + "&fa132=Board%20Game%20Expansions,Board%20Games")

	// Mystery Bay
	s.OnHTML(".s1Dk5D", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "mystery-bay") {
			return
		}

		raw_price := e.ChildText(".s1qgNB")

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".s2OHAT"),
			StoreId:    3,
			StoreThumb: "",
			Stock:      e.ChildAttr(".s1uedD", "aria-disabled") == "false",
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".s2OHAT", "href"),
		})
	})
	q.AddURL("https://www.mystery-bay.com/search-results?q=" + enc)

	// V Games
	s.OnHTML("li.product.type-product", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "v-games") {
			return
		}

		raw_price := e.ChildText(".price")

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".woocommerce-loop-product__title"),
			StoreId:    5,
			StoreThumb: e.ChildAttr(".woocommerce-LoopProduct-link.woocommerce-loop-product__link img", "data-src"),
			Stock:      childHasClass(e, ".button.product_type_simple", "add_to_cart_button"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".woocommerce-LoopProduct-link", "href"),
		})
	})
	q.AddURL("https://store.v-games.gr/?s=" + enc + "&post_type=product&dgwt_wcas=1")

	// kaissa.eu
	s.OnHTML("article.product", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "kaissa.eu") {
			return
		}

		raw_price := e.ChildText(".price")

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".caption"),
			StoreId:    6,
			StoreThumb: e.ChildAttr(".photo a img", "src"),
			Stock:      !childHasClass(e, ".add-to-cart input", "stock-update"),
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".photo a", "href")),
		})
	})
	q.AddURL("https://www.kaissa.eu/products/search?query=" + enc)

	// Meeple Planet
	s.OnHTML("div.product-small.product-type-simple", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "meeple-planet") {
			return
		}

		raw_price := e.ChildText(".amount")

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".name"),
			StoreId:    7,
			StoreThumb: e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			Stock:      !hasClass(e, "out-of-stock"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		})
	})
	q.AddURL("https://meeple-planet.com/?s=" + enc + "&post_type=product")

	// efantasy
	s.OnHTML("div.product.product-box", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "efantasy") {
			return
		}

		raw_price := e.ChildText(".product-price a strong")

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".product-title"),
			StoreId:    8,
			StoreThumb: e.ChildAttr(".product-image a img", "src"),
			Stock:      !(e.ChildAttr(".product-actions button.cartbutton", "data-stock") == "0"),
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".product-title a", "href")),
		})
	})
	q.AddURL("https://www.efantasy.gr/en/products/search=" + enc + "/sort=score/c-31-board-games")

	// kaissagames.com
	s.OnHTML("li.item.product-item", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "kaissagames") {
			return
		}

		raw_price := e.ChildText(".price")

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".name"),
			StoreId:    9,
			StoreThumb: e.ChildAttr(".product-image-photo", "src"),
			Stock:      !childHasClass(e, "div.stock", "unavailable"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".name a", "href"),
		})
	})

	s.OnHTML("div.product-info-main", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "kaissagames") {
			return
		}

		raw_price := e.ChildText(".price")

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".base"),
			StoreId:    9,
			StoreThumb: e.ChildAttr(".MagicZoom figure img", "src"),
			Stock:      !childHasClass(e, "div.new-outofstock", "new-outofstock"),
			Price:      getPrice(raw_price),
			Url:        e.Request.URL.String(),
		})
	})
	q.AddURL("https://kaissagames.com/b2c_gr/catalogsearch/result/?q=" + enc)

	// Meeple on Board
	s.OnHTML("div.product-small.product-type-simple", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "meepleonboard") {
			return
		}

		raw_price := e.ChildText(".amount")

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".name"),
			StoreId:    10,
			StoreThumb: e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			Stock:      !hasClass(e, "out-of-stock"),
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
		})
	})
	q.AddURL("https://meepleonboard.gr/?s=" + enc + "&post_type=product")

	// Games Hero
	s.OnHTML("div.product-layout", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.gameshero.gr") {
			return
		}

		raw_price := e.ChildText(".price-new")

		if raw_price == "" {
			raw_price = e.ChildText(".price-normal")
		}

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".name"),
			StoreId:    11,
			StoreThumb: e.ChildAttr("img.img-responsive", "src"),
			Stock:      !hasClass(e, "out-of-stock"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr("a.product-img", "href"),
		})
	})
	q.AddURL("https://www.gameshero.gr/index.php?route=product/search&search=" + enc + "&description=true")

	// Politeia
	s.OnHTML(".browse-page-block", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), "www.politeianet.gr") {
			return
		}

		raw_price := e.ChildText(".productPrice")

		rs = append(rs, models.Price{
			BoardgameId: models.JsonNullInt64{
				Int64: boardgame.Id,
				Valid: true,
			},
			Name:       e.ChildText(".browse-product-title"),
			StoreId:    12,
			StoreThumb: e.ChildAttr(".browseProductImage", "src"),
			Stock:      childHasClass(e, "input", "addtocart_button_module"),
			Price:      getPrice(raw_price),
			Url:        e.ChildAttr(".browse-product-title", "href"),
		})
	})
	q.AddURL("https://www.politeianet.gr/index.php?option=com_virtuemart&Itemid=89&keyword=" + enc + "&limitstart=0&category_id=943")
	q.Run(s)

	retval := []models.Price{}
	re := regexp.MustCompile(`(?s)\((.*)\)`)

	database, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	for _, price := range rs {
		tmp := re.ReplaceAllString(price.Name, "")
		tmp = strings.ToLower(tmp)
		price.Levenshtein = levenshtein.ComputeDistance(cmd, tmp)
		price.Hamming = Hamming(cmd, tmp)

		if batch_id != nil {
			price.Batch = batch_id.Int64
		}

		if price.Price == 0 {
			continue
		}

		if price.Hamming > 10 {
			continue
		}

		if len(cmd) == price.Hamming {
			continue
		}

		retval = append(retval, price)
	}

	sort.Slice(retval, func(i int, j int) bool {
		return retval[i].Levenshtein < retval[j].Levenshtein
	})

	if batch_id != nil {
		for _, price := range retval {
			_, err := create(database, price)
			if err != nil {
				return nil, err
			}
		}
	}

	err = storage.Clear()
	if err != nil {
		return nil, err
	}

	return retval, nil
}
