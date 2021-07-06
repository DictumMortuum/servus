package scraper

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type Data struct{}

func (obj Data) GetTable() string {
	return "tboardgamescraperdata"
}

func getData(db *sqlx.DB, id int64) (*models.ScraperData, error) {
	var rs models.ScraperData

	err := db.QueryRowx(`select * from tboardgamescraperdata where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func (obj Data) Get(db *sqlx.DB, id int64) (interface{}, error) {
	return getData(db, id)
}

func (obj Data) GetList(db *sqlx.DB, query string, args ...interface{}) (interface{}, error) {
	var rs []models.ScraperData

	err := db.Select(&rs, query, args...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Data) Create(db *sqlx.DB, query string, data map[string]interface{}) (interface{}, error) {
	var game models.ScraperData

	if val, ok := data["store_id"]; ok {
		game.StoreId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	} else {
		return nil, errors.New("please provide a 'store_id' parameter")
	}

	if val, ok := data["boardgame_id"]; ok {
		game.BoardgameId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	} else {
		game.BoardgameId = models.JsonNullInt64{
			Int64: -1,
			Valid: false,
		}
		//return nil, errors.New("please provide a 'boardgame_id' parameter")
	}

	if val, ok := data["title"]; ok {
		game.Title = val.(string)
	} else {
		return nil, errors.New("please provide a 'title' parameter")
	}

	if val, ok := data["link"]; ok {
		game.Link = val.(string)
	} else {
		return nil, errors.New("please provide a 'link' parameter")
	}

	if val, ok := data["sku"]; ok {
		game.SKU = val.(string)
	} else {
		game.SKU = ""
	}

	if val, ok := data["active"]; ok {
		t, err := time.Parse("2006-01-02T15:04:05-0700", val.(string))
		if err != nil {
			return nil, err
		}

		game.Active = t
	}

	game.CrDate = time.Now()

	rs, err := db.NamedExec(query, &game)
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}

	game.Id = id
	return game, nil
}

func (obj Data) Update(db *sqlx.DB, query string, id int64, data map[string]interface{}) (interface{}, error) {
	game, err := getData(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := data["store_id"]; ok {
		game.StoreId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	}

	if val, ok := data["boardgame_id"]; ok {
		game.BoardgameId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	}

	if val, ok := data["title"]; ok {
		game.Title = val.(string)
	}

	if val, ok := data["link"]; ok {
		game.Link = val.(string)
	}

	if val, ok := data["sku"]; ok {
		game.SKU = val.(string)
	}

	_, err = db.NamedExec(query, &game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (obj Data) Delete(db *sqlx.DB, query string, id int64) (interface{}, error) {
	rs, err := getData(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
