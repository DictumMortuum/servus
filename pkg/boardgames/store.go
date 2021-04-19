package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func GetStore(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getStore(db, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetStoreList(db *sqlx.DB, args models.Args) (interface{}, int, error) {
	var rs []models.Store

	sql, err := args.List(`
	select
		*
	from
		tboardgamestore
	`)
	if err != nil {
		return nil, -1, err
	}

	if len(args.Id) > 0 {
		query, ids, err := sqlx.In(sql.String(), args.Id)
		if err != nil {
			return nil, -1, err
		}

		err = db.Select(&rs, db.Rebind(query), ids...)
		if err != nil {
			return nil, -1, err
		}
	} else {
		err = db.Select(&rs, sql.String())
		if err != nil {
			return nil, -1, err
		}
	}

	return rs, len(rs), nil
}

func CreateStore(db *sqlx.DB, data map[string]interface{}) (interface{}, error) {
	var player models.Store

	if val, ok := data["name"]; ok {
		player.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	sql := `
	insert into tboardgamestore (
		name
	) values (
		:name
	)`

	rs, err := db.NamedExec(sql, &player)
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}

	player.Id = id

	return player, nil
}

func UpdateStore(db *sqlx.DB, id int64, data map[string]interface{}) (interface{}, error) {
	player, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := data["name"]; ok {
		player.Name = val.(string)
	}

	sql := `
	update
		tboardgamestore
	set
		name = :name
	where
		id = :id`

	_, err = db.NamedExec(sql, &player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func DeleteStore(db *sqlx.DB, id int64) (interface{}, error) {
	player, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	sql := `
	delete from
		tboardgamestore
	where
		id = :id`

	_, err = db.NamedExec(sql, &player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func getStore(db *sqlx.DB, id int64) (*models.Store, error) {
	var retval models.Store

	err := db.QueryRowx(`select * from tboardgamestore where id = ?`, id).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}
