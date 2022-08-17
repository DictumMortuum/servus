package boardgames

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func getMapping(db *sqlx.DB, id int64) (*models.Mapping, error) {
	var rs models.Mapping

	err := db.QueryRowx(`select * from tboardgamepricesmap where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func GetMapping(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getMapping(db, args.Id)
}

func GetListMapping(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Mapping

	var count []int
	err := db.Select(&count, "select 1 from tboardgamepricesmap")
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", len(count)))

	sql, err := args.List(`
		select * from tboardgamepricesmap
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

func CreateMapping(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs models.Mapping

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		err := fn(args.Data, true)
		if err != nil {
			return nil, err
		}
	}

	query, err := args.Insert("tboardgamepricesmap")
	if err != nil {
		return nil, err
	}

	row, err := db.NamedExec(query.String(), &rs)
	if err != nil {
		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	rs.Id = id
	return rs, nil
}

func UpdateMapping(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getMapping(db, args.Id)
	if err != nil {
		return nil, err
	}

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		fn(args.Data, false)
	}

	sql, err := args.Update("tboardgamepricesmap")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func DeleteMapping(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getMapping(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgamepricesmap where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
