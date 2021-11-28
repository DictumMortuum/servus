package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Store struct{}

func getStore(db *sqlx.DB, id int64) (*models.Store, error) {
	var rs models.Store

	err := db.QueryRowx(`select * from tboardgamestores where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func (obj Store) Get(db *sqlx.DB, id int64) (interface{}, error) {
	return getStore(db, id)
}

func (obj Store) GetList(db *sqlx.DB, args models.QueryBuilder) (interface{}, int, error) {
	var rs []models.Store

	var count []int
	err := db.Select(&count, "select 1 from tboardgamestores")
	if err != nil {
		return nil, -1, err
	}

	sql, err := args.List(`
		select * from tboardgamestores
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

func (obj Store) Create(db *sqlx.DB, qb models.QueryBuilder) (interface{}, error) {
	var store models.Store

	if val, ok := qb.Data["name"]; ok {
		store.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	query, err := qb.Insert("tboardgamestores")
	if err != nil {
		return nil, err
	}

	rs, err := db.NamedExec(query.String(), &store)
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}

	store.Id = id
	return store, nil
}

func (obj Store) Update(db *sqlx.DB, id int64, qb models.QueryBuilder) (interface{}, error) {
	rs, err := getStore(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := qb.Data["name"]; ok {
		rs.Name = val.(string)
	}

	sql, err := qb.Update("tboardgamestores")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Store) Delete(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getStore(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgamestores where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
