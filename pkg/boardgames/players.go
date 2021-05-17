package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Player struct{}

func (obj Player) GetTable() string {
	return "tboardgameplayers"
}

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

func (obj Player) GetList(db *sqlx.DB, query string, args ...interface{}) (interface{}, error) {
	var rs []models.Boardgame

	err := db.Select(&rs, query, args...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Player) Create(db *sqlx.DB, query string, data map[string]interface{}) (interface{}, error) {
	var rs models.Player

	if val, ok := data["name"]; ok {
		rs.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	_, err := db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Player) Update(db *sqlx.DB, query string, id int64, data map[string]interface{}) (interface{}, error) {
	rs, err := getPlayer(db, id)
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

func (obj Player) Delete(db *sqlx.DB, query string, id int64) (interface{}, error) {
	rs, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
