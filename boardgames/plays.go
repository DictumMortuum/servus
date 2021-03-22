package boardgames

import (
	"bytes"
	"errors"
	"github.com/DictumMortuum/servus/generic"
	"github.com/jmoiron/sqlx"
	"strconv"
	"text/template"
	"time"
)

type PlayModel struct {
	Id          int64     `db:"id" json:"id"`
	CrDate      time.Time `db:"cr_date" json:"cr_date"`
	Date        time.Time `db:"date" json:"date"`
	BoardgameId int64     `db:"boardgame_id"`
}

func GetPlay(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getPlay(db, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetPlayList(db *sqlx.DB, args generic.Args) (interface{}, int, error) {
	var rs []PlayModel

	sql := `
	select
		*
	from
		tboardgameplays
	{{ if gt (len .Id) 0 }}
	where
		{{ .RefKey }} in (?)
	{{ else if gt (len .FilterVal) 0 }}
	where
		{{ .FilterKey }} = "{{ .FilterVal }}"
	{{ end }}
	{{ if gt (len .Sort) 0 }}
	order by {{ .Sort }} {{ .Order }}
	{{ end }}
	{{ if eq (len .Range) 2 }}
	limit {{ index .Range 0 }}, {{ .Page }}
	{{ end }}`

	var tpl bytes.Buffer
	t := template.Must(template.New("playslist").Parse(sql))
	err := t.Execute(&tpl, args)
	if err != nil {
		return nil, -1, err
	}

	if len(args.Id) > 0 {
		query, ids, err := sqlx.In(tpl.String(), args.Id)
		if err != nil {
			return nil, -1, err
		}

		err = db.Select(&rs, db.Rebind(query), ids...)
		if err != nil {
			return nil, -1, err
		}
	} else {
		err = db.Select(&rs, tpl.String())
		if err != nil {
			return nil, -1, err
		}
	}

	return rs, len(rs), nil
}

func CreatePlay(db *sqlx.DB, data map[string]interface{}) (interface{}, error) {
	var play PlayModel
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

func getPlay(db *sqlx.DB, id int64) (*PlayModel, error) {
	var rs PlayModel

	err := db.QueryRowx(`select * from tboardgameplays where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}
