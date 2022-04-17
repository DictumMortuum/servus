package bgg

import (
	"database/sql"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"strings"
	"time"
)

func FetchTopArt(col *gin.Context) {
	data := []int64{}

	database, err := db.Conn()
	if err != nil {
		util.Error(col, err)
		return
	}
	defer database.Close()

	rs, err := getTopBoardgames(database)
	if err != nil {
		util.Error(col, err)
		return
	}

	for _, item := range rs {
		_, err := FetchBoardgame(database, item.Id)
		if err != nil {
			util.Error(col, err)
			return
		}
		data = append(data, item.Id)
		time.Sleep(3 * time.Second)
	}

	util.Success(col, &data)
}

func GetTopBoardgames(col *gin.Context) {
	data := []map[string]interface{}{}

	database, err := db.Conn()
	if err != nil {
		util.Error(col, err)
		return
	}
	defer database.Close()

	c := colly.NewCollector(
		colly.AllowedDomains("boardgamegeek.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*boardgamegeek.com.*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	c.OnHTML("#collectionitems tbody tr", func(e *colly.HTMLElement) {
		raw_rank := e.ChildText(".collection_rank")
		name := e.ChildText(".collection_objectname div a")
		url := e.ChildAttr(".collection_objectname div a", "href")
		// thumb := e.ChildAttr(".collection_thumbnail a img", "src")
		tokens := strings.Split(url, "/")

		if len(tokens) == 4 {
			raw_id, _ := strconv.ParseInt(tokens[2], 10, 64)
			rank, _ := strconv.ParseInt(raw_rank, 10, 64)

			d := map[string]interface{}{
				"name": name,
				"rank": rank,
				"url":  url,
				"id":   raw_id,
			}

			id, err := exists(database, d)
			if err != nil {
				util.Error(col, err)
				return
			}

			if id == nil {
				_, err := create(database, d)
				if err != nil {
					util.Error(col, err)
					return
				}
			} else {
				_, err := updateRank(database, d)
				if err != nil {
					util.Error(col, err)
					return
				}
			}
		}
	})

	for i := 1; i <= 1335; i++ {
		c.Visit(fmt.Sprintf("https://boardgamegeek.com/browse/boardgame/page/%d", i))
	}

	util.Success(col, &data)
}

func exists(db *sqlx.DB, payload map[string]interface{}) (*models.JsonNullInt64, error) {
	var id models.JsonNullInt64

	q := `select id from tboardgames where id = :id`
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

func create(db *sqlx.DB, payload map[string]interface{}) (bool, error) {
	q := `insert into tboardgames (id,name,rank,configured) values (:id,:name,:rank,0)`

	rs, err := db.NamedExec(q, payload)
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	log.Printf("Boardgame [%s] with id [%d] created: [%d]\n", payload["name"], payload["id"], rows)
	return rows > 0, nil
}

func updateRank(db *sqlx.DB, payload map[string]interface{}) (bool, error) {
	q := `
		update
			tboardgames
		set
			rank = :rank
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

func update(db *sqlx.DB, payload map[string]interface{}) (bool, error) {
	q := `
		update
			tboardgames
		set
			rank = :rank,
			thumb = :thumb
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

	log.Printf("Boardgame [%s] with id [%d] updated: [%d]\n", payload["name"], payload["id"], rows)
	return rows > 0, nil
}

func getTopBoardgames(db *sqlx.DB) ([]models.Boardgame, error) {
	var rs []models.Boardgame

	sql := `
		select
			*
		from
			tboardgames
		where
			rank <= 300 and rank != 0
		order by rank
	`

	err := db.Select(&rs, sql)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
