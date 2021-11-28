package scraper

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type Data struct{}

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

func (obj Data) GetList(db *sqlx.DB, args models.QueryBuilder) (interface{}, int, error) {
	var rs []models.ScraperData

	var count []int
	err := db.Select(&count, "select 1 from tboardgamescraperdata")
	if err != nil {
		return nil, -1, err
	}

	sql, err := args.List(`
		select * from tboardgamescraperdata
	`)
	if err != nil {
		return nil, -1, err
	}

	query, ids, err := sqlx.In(sql.String(), args.Id)
	if err != nil {
		query = sql.String()
	} else {
		query = db.Rebind(query)
	}

	err = db.Select(&rs, query, ids...)
	if err != nil {
		return nil, -1, err
	}

	return rs, len(count), nil
}

func (obj Data) Create(db *sqlx.DB, qb models.QueryBuilder) (interface{}, error) {
	var game models.ScraperData

	if val, ok := qb.Data["store_id"]; ok {
		game.StoreId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	} else {
		return nil, errors.New("please provide a 'store_id' parameter")
	}

	if val, ok := qb.Data["boardgame_id"]; ok {
		game.BoardgameId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	} else {
		game.BoardgameId = models.JsonNullInt64{
			Int64: -1,
			Valid: false,
		}
	}

	if val, ok := qb.Data["title"]; ok {
		game.Title = val.(string)
	} else {
		return nil, errors.New("please provide a 'title' parameter")
	}

	if val, ok := qb.Data["link"]; ok {
		game.Link = val.(string)
	} else {
		return nil, errors.New("please provide a 'link' parameter")
	}

	if val, ok := qb.Data["sku"]; ok {
		game.SKU = val.(string)
	} else {
		game.SKU = ""
	}

	if val, ok := qb.Data["active"]; ok {
		t, err := time.Parse("2006-01-02T15:04:05-0700", val.(string))
		if err != nil {
			return nil, err
		}

		game.Active = t
	}

	game.CrDate = time.Now()

	query, err := qb.Insert("tboardgamescraperdata")
	if err != nil {
		return nil, err
	}

	rs, err := db.NamedExec(query.String(), &game)
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

func (obj Data) Update(db *sqlx.DB, id int64, qb models.QueryBuilder) (interface{}, error) {
	game, err := getData(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := qb.Data["store_id"]; ok {
		game.StoreId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	}

	if val, ok := qb.Data["boardgame_id"]; ok {
		game.BoardgameId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	}

	if val, ok := qb.Data["title"]; ok {
		game.Title = val.(string)
	}

	if val, ok := qb.Data["link"]; ok {
		game.Link = val.(string)
	}

	if val, ok := qb.Data["sku"]; ok {
		game.SKU = val.(string)
	}

	sql, err := qb.Update("tboardgamescraperdata")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (obj Data) Delete(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getData(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgamescraperdata where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
