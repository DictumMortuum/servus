package bgg

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
	"strconv"
)

func boardgameExists(db *sqlx.DB, payload map[string]interface{}) (*models.JsonNullInt64, error) {
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

func FetchBoardgameIfNotExists(db *sqlx.DB, id models.JsonNullInt64) (*models.Boardgame, error) {
	if id.Valid {
		exists, err := boardgameExists(db, map[string]interface{}{
			"id": id.Int64,
		})
		if err != nil {
			return nil, err
		}

		if exists == nil {
			_, err := FetchBoardgame(db, id.Int64)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func FetchBoardgame(db *sqlx.DB, id int64) (*models.Boardgame, error) {
	var rs models.Boardgame

	link := fmt.Sprintf("https://www.boardgamegeek.com/xmlapi2/thing?id=%d&stats=1", id)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	conn := &http.Client{}
	resp, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tmp := BggThing{}
	err = xml.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}

	for _, item := range tmp.Items {
		d := map[string]interface{}{
			"name":    item.Name.Value,
			"rank":    getRank(item.Statistics.Ratings.Ranks.Ranks),
			"id":      id,
			"preview": item.Image,
		}

		e, err := exists(db, d)
		if err != nil {
			return nil, err
		}

		if e == nil {
			_, err := create(db, d)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := update(db, d)
			if err != nil {
				return nil, err
			}
		}
	}

	return &rs, nil
}

func getRank(ranks []Rank) int64 {
	for _, rank := range ranks {
		if rank.Name == "Board Game Rank" {
			v, err := strconv.ParseInt(rank.Value, 10, 64)
			if err != nil {
				return 0
			}
			return v
		}
	}

	return 0
}
