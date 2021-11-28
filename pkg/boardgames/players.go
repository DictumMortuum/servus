package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Player struct{}

func getPlayer(db *sqlx.DB, id int64) (*models.Player, error) {
	var rs models.Player

	err := db.QueryRowx(`select * from tboardgameplayers where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func (obj Player) Get(db *sqlx.DB, id int64) (interface{}, error) {
	return getPlayer(db, id)
}

func (obj Player) GetList(db *sqlx.DB, args models.QueryBuilder) (interface{}, int, error) {
	var rs []models.Player

	var count []int
	err := db.Select(&count, "select 1 from tboardgameplayers")
	if err != nil {
		return nil, -1, err
	}

	sql, err := args.List(`
		select * from tboardgameplayers
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

func (obj Player) Create(db *sqlx.DB, qb models.QueryBuilder) (interface{}, error) {
	var rs models.Player

	if val, ok := qb.Data["name"]; ok {
		rs.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	query, err := qb.Insert("tboardgameplayers")
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

func (obj Player) Update(db *sqlx.DB, id int64, qb models.QueryBuilder) (interface{}, error) {
	rs, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	sql, err := qb.Update("tboardgameplayers")
	if err != nil {
		return nil, err
	}

	if val, ok := qb.Data["name"]; ok {
		rs.Name = val.(string)
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Player) Delete(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgameplayers where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
