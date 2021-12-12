package boardgames

import (
	"errors"
	"fmt"
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

func (obj Boardgame) Get(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getBoardgame(db, args.Id)
}

func (obj Boardgame) GetList(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Boardgame

	var count []int
	err := db.Select(&count, "select 1 from tboardgames")
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", len(count)))

	sql, err := args.List(`
		select * from tboardgames
	`)
	if err != nil {
		return nil, err
	}

	query, ids, err := sqlx.In(sql.String(), args.Ids)
	if err != nil {
		query = sql.String()
	} else {
		query = db.Rebind(query)
	}

	err = db.Select(&rs, query, ids...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Boardgame) Create(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs models.Boardgame

	if val, ok := args.Data["id"]; ok {
		rs.Id = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	if val, ok := args.Data["name"]; ok {
		rs.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	if val, ok := args.Data["data"]; ok {
		err := rs.Data.Scan(val)
		if err != nil {
			return nil, err
		}
	} else {
		rs.Data = nil
	}

	query, err := args.Insert("tboardgames")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(query.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Boardgame) Update(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getBoardgame(db, args.Id)
	if err != nil {
		return nil, err
	}

	if val, ok := args.Data["name"]; ok {
		rs.Name = val.(string)
	}

	if val, ok := args.Data["data"]; ok {
		err := rs.Data.Scan(val)
		if err != nil {
			return nil, err
		}
	}

	sql, err := args.Update("tboardgames")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Boardgame) Delete(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getBoardgame(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgames where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
