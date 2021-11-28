package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type Boardgame struct{}

func getBoardgame(db *sqlx.DB, id int64) (*models.Boardgame, error) {
	var rs models.Boardgame

	err := db.QueryRowx(`select * from tboardgames where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	rs.Date = time.Now().AddDate(1, 0, 0)

	return &rs, nil
}

func (obj Boardgame) Get(db *sqlx.DB, id int64) (interface{}, error) {
	return getBoardgame(db, id)
}

func (obj Boardgame) GetList(db *sqlx.DB, args models.QueryBuilder) (interface{}, int, error) {
	var rs []models.Boardgame

	var count []int
	err := db.Select(&count, "select 1 from tboardgames")
	if err != nil {
		return nil, -1, err
	}

	sql, err := args.List(`
		select * from tboardgames
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

func (obj Boardgame) Create(db *sqlx.DB, qb models.QueryBuilder) (interface{}, error) {
	var rs models.Boardgame

	if val, ok := qb.Data["id"]; ok {
		rs.Id = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	if val, ok := qb.Data["name"]; ok {
		rs.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	query, err := qb.Insert("tboardgames")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(query.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Boardgame) Update(db *sqlx.DB, id int64, qb models.QueryBuilder) (interface{}, error) {
	rs, err := getBoardgame(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := qb.Data["name"]; ok {
		rs.Name = val.(string)
	}

	sql, err := qb.Update("tboardgames")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Boardgame) Delete(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getBoardgame(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgames where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
