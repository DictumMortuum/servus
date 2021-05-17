package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type Play struct{}

func (obj Play) GetTable() string {
	return "tboardgameplays"
}

func getPlay(db *sqlx.DB, id int64) (*models.Play, error) {
	var rs models.Play

	err := db.QueryRowx(`select * from tboardgameplays where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func (obj Play) Get(db *sqlx.DB, id int64) (interface{}, error) {
	return getPlay(db, id)
}

func (obj Play) GetList(db *sqlx.DB, query string, args ...interface{}) (interface{}, error) {
	var rs []models.Play

	err := db.Select(&rs, query, args...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Play) Create(db *sqlx.DB, query string, data map[string]interface{}) (interface{}, error) {
	var rs models.Play

	if val, ok := data["date"]; ok {
		t, err := time.Parse("2006-01-02T15:04:05-0700", val.(string))
		if err != nil {
			return nil, err
		}

		rs.Date = t
	}

	if val, ok := data["boardgame_id"]; ok {
		rs.Id = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'boardgame_id' parameter")
	}

	rs.CrDate = time.Now()

	_, err := db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Play) Update(db *sqlx.DB, query string, id int64, data map[string]interface{}) (interface{}, error) {
	rs, err := getPlay(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := data["date"]; ok {
		t, err := time.Parse("2006-01-02T15:04:05-0700", val.(string))
		if err != nil {
			return nil, err
		}

		rs.Date = t
	}

	_, err = db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (obj Play) Delete(db *sqlx.DB, query string, id int64) (interface{}, error) {
	rs, err := getPlay(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(query, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
