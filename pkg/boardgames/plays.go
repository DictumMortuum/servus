package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

func GetPlay(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getPlay(db, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetPlayList(db *sqlx.DB, args models.Args) (interface{}, int, error) {
	var rs []models.Play

	sql, err := args.List(`
	select
		*
	from
		tboardgameplays
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

func CreatePlay(db *sqlx.DB, data map[string]interface{}) (interface{}, error) {
	var play models.Play
	var err error

	if val, ok := data["cr_date"]; ok {
		play.Date, err = time.Parse(time.RFC3339, val.(string))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("please provide a 'cr_date' parameter")
	}

	if val, ok := data["date"]; ok {
		play.CrDate, err = time.Parse(time.RFC3339, val.(string))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("please provide a 'date' parameter")
	}

	if val, ok := data["boardgame_id"]; ok {
		play.BoardgameId, err = strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("please provide a 'boardgame_id' parameter")
	}

	sql := `
	insert into tboardgameplays (
		boardgame_id,
		cr_date,
		date
	) values (
		:boardgame_id,
		:cr_date,
		:date
	)`

	rs, err := db.NamedExec(sql, &play)
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}

	play.Id = id

	return play, nil
}

//TODO
func UpdatePlay(db *sqlx.DB, id int64, data map[string]interface{}) (interface{}, error) {
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

//TODO
func DeletePlay(db *sqlx.DB, id int64) (interface{}, error) {
	play, err := getPlay(db, id)
	if err != nil {
		return nil, err
	}

	sql := `
	delete from
		tboardgameplayers
	where
		id = :id`

	_, err = db.NamedExec(sql, &play)
	if err != nil {
		return nil, err
	}

	return play, nil
}

func getPlay(db *sqlx.DB, id int64) (*models.Play, error) {
	var rs models.Play

	err := db.QueryRowx(`select * from tboardgameplays where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}
