package boardgames

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func getPlayer(db *sqlx.DB, id int64) (*models.Player, error) {
	var rs models.Player

	err := db.QueryRowx(`select * from tboardgameplayers where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func GetPlayer(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getPlayer(db, args.Id)
}

func GetListPlayer(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Player

	var count []int
	err := db.Select(&count, "select 1 from tboardgameplayers")
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", len(count)))

	sql, err := args.List(`
		select * from tboardgameplayers
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

func CreatePlayer(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs models.Player

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		err := fn(args.Data, true)
		if err != nil {
			return nil, err
		}
	}

	query, err := args.Insert("tboardgameplayers")
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

func UpdatePlayer(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getPlayer(db, args.Id)
	if err != nil {
		return nil, err
	}

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		fn(args.Data, false)
	}

	sql, err := args.Update("tboardgameplayers")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func DeletePlayer(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getPlayer(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgameplayers where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
