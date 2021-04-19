package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func GetPlayer(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetPlayerList(db *sqlx.DB, args models.Args) (interface{}, int, error) {
	var rs []models.Player

	sql, err := args.List(`
	select
		*
	from
		tboardgameplayers
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

func CreatePlayer(db *sqlx.DB, data map[string]interface{}) (interface{}, error) {
	var player models.Player

	if val, ok := data["name"]; ok {
		player.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	sql := `
	insert into tboardgameplayers (
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

func UpdatePlayer(db *sqlx.DB, id int64, data map[string]interface{}) (interface{}, error) {
	player, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := data["name"]; ok {
		player.Name = val.(string)
	}

	sql := `
	update
		tboardgameplayers
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

func DeletePlayer(db *sqlx.DB, id int64) (interface{}, error) {
	player, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	sql := `
	delete from
		tboardgameplayers
	where
		id = :id`

	_, err = db.NamedExec(sql, &player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func createPlayer(db *sqlx.DB, data models.Player) error {
	sql := `
	insert into tboardgameplayers (
		name
	) values (
		:name
	)`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func getPlayer(db *sqlx.DB, id int64) (*models.Player, error) {
	var retval models.Player

	err := db.QueryRowx(`select * from tboardgameplayers where id = ?`, id).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}

func getPlayerByName(db *sqlx.DB, name string) (*models.Player, error) {
	var retval models.Player

	err := db.QueryRowx(`select * from tboardgameplayers where name = ?`, name).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}
