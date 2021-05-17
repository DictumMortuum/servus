package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Stats struct{}

func (obj Stats) GetTable() string {
	return "tboardgamestats"
}

func getStats(db *sqlx.DB, id int64) (*models.Stats, error) {
	var rs models.Stats

	err := db.QueryRowx(`select * from tboardgamestats where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func (obj Stats) Get(db *sqlx.DB, id int64) (interface{}, error) {
	return getStats(db, id)
}

func (obj Stats) GetList(db *sqlx.DB, query string, args ...interface{}) (interface{}, error) {
	var rs []models.Stats

	err := db.Select(&rs, query, args...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Stats) Create(db *sqlx.DB, query string, data map[string]interface{}) (interface{}, error) {
	var stats models.Stats

	if val, ok := data["play_id"]; ok {
		stats.PlayId = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'play_id' parameter")
	}

	if val, ok := data["player_id"]; ok {
		stats.PlayerId = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'player_id' parameter")
	}

	if val, ok := data["json"]; ok {
		stats.Data.Scan(val)
	} else {
		return nil, errors.New("please provide a 'json' parameter")
	}

	rs, err := db.NamedExec(query, &stats)
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}

	stats.Id = id
	return stats, nil
}

func (obj Stats) Update(db *sqlx.DB, query string, id int64, data map[string]interface{}) (interface{}, error) {
	rs, err := getStats(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := data["play_id"]; ok {
		rs.PlayId = int64(val.(float64))
	}

	if val, ok := data["player_id"]; ok {
		rs.PlayerId = int64(val.(float64))
	}

	if val, ok := data["json"]; ok {
		rs.Data.Scan(val)
	}

	_, err = db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Stats) Delete(db *sqlx.DB, query string, id int64) (interface{}, error) {
	rs, err := getStats(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
