package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type Boardgame struct{}

func (obj Boardgame) GetTable() string {
	return "tboardgames"
}

func getBoardgame(db *sqlx.DB, id int64) (*models.Boardgame, error) {
	var rs models.Boardgame

	err := db.QueryRowx(`select * from tboardgames where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	rs.Date = time.Now().AddDate(1, 0, 0)

	return &rs, nil
}

func (obj Boardgame) Get(db *sqlx.DB, id int64) (interface{}, error) {
	return getBoardgame(db, id)
}

func (obj Boardgame) GetList(db *sqlx.DB, query string, args ...interface{}) (interface{}, error) {
	var rs []models.Boardgame

	err := db.Select(&rs, query, args...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Boardgame) Create(db *sqlx.DB, query string, data map[string]interface{}) (interface{}, error) {
	var rs models.Boardgame

	if val, ok := data["id"]; ok {
		rs.Id = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

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

func (obj Boardgame) Update(db *sqlx.DB, query string, id int64, data map[string]interface{}) (interface{}, error) {
	rs, err := getBoardgame(db, id)
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

func (obj Boardgame) Delete(db *sqlx.DB, query string, id int64) (interface{}, error) {
	rs, err := getBoardgame(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
