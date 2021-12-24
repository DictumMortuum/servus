package boardgames

import (
	"errors"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func getStore(db *sqlx.DB, id int64) (*models.Store, error) {
	var rs models.Store

	err := db.QueryRowx(`select * from tboardgamestores where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func GetStore(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getStore(db, args.Id)
}

func GetListStore(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Store

	var count []int
	err := db.Select(&count, "select 1 from tboardgamestores")
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", len(count)))

	sql, err := args.List(`
		select * from tboardgamestores
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

func CreateStore(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var store models.Store

	if val, ok := args.Data["name"]; ok {
		store.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	query, err := args.Insert("tboardgamestores")
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

func UpdateStore(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getStore(db, args.Id)
	if err != nil {
		return nil, err
	}

	if val, ok := args.Data["name"]; ok {
		rs.Name = val.(string)
	}

	sql, err := args.Update("tboardgamestores")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func DeleteStore(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getStore(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgamestores where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
