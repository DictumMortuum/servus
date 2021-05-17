package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Store struct{}

func (obj Store) GetTable() string {
	return "tboardgamestores"
}

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

func (obj Store) GetList(db *sqlx.DB, query string, args ...interface{}) (interface{}, error) {
	var rs []models.Store

	err := db.Select(&rs, query, args...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Store) Create(db *sqlx.DB, query string, data map[string]interface{}) (interface{}, error) {
	var store models.Store

	if val, ok := data["name"]; ok {
		store.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	rs, err := db.NamedExec(query, &store)
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

func (obj Store) Update(db *sqlx.DB, query string, id int64, data map[string]interface{}) (interface{}, error) {
	rs, err := getStore(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := data["name"]; ok {
		rs.Name = val.(string)
	}

	_, err = db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Store) Delete(db *sqlx.DB, query string, id int64) (interface{}, error) {
	rs, err := getStore(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
