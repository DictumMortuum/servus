package boardgames

import (
	"bytes"
	"errors"
	"github.com/DictumMortuum/servus/generic"
	"github.com/jmoiron/sqlx"
	"text/template"
)

type BoardgameModel struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

func GetBoardgame(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getBoardgame(db, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetBoardgameList(db *sqlx.DB, args generic.Args) (interface{}, int, error) {
	var rs []BoardgameModel

	sql := `
	select
		*
	from
		tboardgames
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
	t := template.Must(template.New("boardgamelist").Parse(sql))
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

func CreateBoardgame(db *sqlx.DB, data map[string]interface{}) (interface{}, error) {
	var player BoardgameModel

	if val, ok := data["name"]; ok {
		player.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	sql := `
	insert into tboardgames (
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

func UpdateBoardgame(db *sqlx.DB, id int64, data map[string]interface{}) (interface{}, error) {
	player, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := data["name"]; ok {
		player.Name = val.(string)
	}

	sql := `
	update
		tboardgames
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

func DeleteBoardgame(db *sqlx.DB, id int64) (interface{}, error) {
	player, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	sql := `
	delete from
		tboardgames
	where
		id = :id`

	_, err = db.NamedExec(sql, &player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func getBoardgame(db *sqlx.DB, id int64) (*BoardgameModel, error) {
	var retval BoardgameModel

	err := db.QueryRowx(`select * from tboardgames where id = ?`, id).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}
