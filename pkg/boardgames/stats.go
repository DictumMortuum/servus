package boardgames

import (
	"bytes"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"text/template"
)

func getStats(db *sqlx.DB, id int64) (*models.Stats, error) {
	var rs models.Stats

	err := db.QueryRowx(`
		select
			s.*,
			pl.name,
			pl.surname
		from
			tboardgamestats s,
			tboardgameplayers pl
		where
			s.player_id = pl.id and
			s.id = ?
	`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func getPlayStats(db *sqlx.DB, id int64) ([]models.Stats, error) {
	var rs []models.Stats

	err := db.Select(&rs, `
		select
			s.*,
			pl.name,
			pl.surname
		from
			tboardgamestats s,
			tboardgameplayers pl
		where
			s.player_id = pl.id and
			s.play_id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetStats(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getStats(db, args.Id)
}

func GetListStats(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []models.Stats{}

	var count []int
	err := db.Select(&count, "select 1 from tboardgamestats")
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", len(count)))

	sql := `
		select
			s.*,
			pl.name,
			pl.surname
		from
			tboardgamestats s,
			tboardgameplayers pl
		where
			s.player_id = pl.id
		{{ if gt (len .Ids) 0 }}
			and s.{{ .RefKey }} in (?)
		{{ else if gt (len .FilterVal) 0 }}
			and s.{{ .FilterKey }} = "{{ .FilterVal }}"
		{{ end }}
		{{ if gt (len .Sort) 0 }}
		order by s.{{ .Sort }} {{ .Order }}
		{{ end }}
		{{ if eq (len .Range) 2 }}
		limit {{ index .Range 0 }}, {{ .Page }}
		{{ end }}`

	var tpl bytes.Buffer
	t := template.Must(template.New("list").Parse(sql))
	err = t.Execute(&tpl, args)
	if err != nil {
		return nil, err
	}

	query, ids, err := sqlx.In(tpl.String(), args.Ids)
	if err != nil {
		query = tpl.String()
	} else {
		query = db.Rebind(query)
	}

	err = db.Select(&rs, query, ids...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func CreateStats(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs models.Stats

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		err := fn(args.Data, true)
		if err != nil {
			return nil, err
		}
	}

	query, err := args.Insert("tboardgamestats")
	if err != nil {
		return nil, err
	}

	stats, err := db.NamedExec(query.String(), &rs)
	if err != nil {
		return nil, err
	}

	id, err := stats.LastInsertId()
	if err != nil {
		return nil, err
	}

	rs.Id = id
	return rs, nil
}

func UpdateStats(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getStats(db, args.Id)
	if err != nil {
		return nil, err
	}

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		fn(args.Data, false)
	}

	sql, err := args.Update("tboardgamestats")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func DeleteStats(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getStats(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgamestats where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
