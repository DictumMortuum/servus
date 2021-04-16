package boardgames

import (
	"bytes"
	"errors"
	"github.com/DictumMortuum/servus/pkg/generic"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"text/template"
)

func GetPlayer(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetPlayerList(db *sqlx.DB, args generic.Args) (interface{}, int, error) {
	var rs []models.Player

	sql := `
	select
		*
	from
		tboardgameplayers
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
	t := template.Must(template.New("playerslist").Parse(sql))
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
