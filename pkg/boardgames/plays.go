package boardgames

import (
	"bytes"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/boardgames/score"
	"github.com/DictumMortuum/servus/pkg/boardgames/trueskill"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"text/template"
)

func getPlay(db *sqlx.DB, id int64) (*models.Play, error) {
	var rs models.Play

	err := db.QueryRowx(`
		select
			p.*,
			g.name,
			g.data
		from
			tboardgameplays p,
			tboardgames g
		where
			p.boardgame_id = g.id and
			p.id = ?
	`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	stats, err := getPlayStats(db, id)
	if err != nil {
		return nil, err
	}
	rs.Stats = stats

	play, err := score.Calculate(db, rs)
	if err != nil {
		return nil, err
	}

	return play, nil
}

func GetPlay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getPlay(db, args.Id)
}

func GetListPlay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []models.Play{}

	var count []int
	err := db.Select(&count, "select 1 from tboardgameplays")
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", len(count)))

	sql := `
		select
			p.*,
			g.name,
			g.data
		from
			tboardgameplays p,
			tboardgames g
		where
			p.boardgame_id = g.id
		{{ if gt (len .Ids) 0 }}
			and p.{{ .RefKey }} in (?)
		{{ else if gt (len .FilterVal) 0 }}
			and p.{{ .FilterKey }} = "{{ .FilterVal }}"
		{{ end }}
		{{ if gt (len .Sort) 0 }}
		order by p.{{ .Sort }} {{ .Order }}
		{{ else }}
		order by p.date asc, p.id
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

	if args.Resources["stats"] {
		retval := []models.Play{}

		for _, item := range rs {
			play, err := score.Calculate(db, item)
			if err != nil {
				return nil, err
			}

			if play != nil {
				retval = append(retval, *play)
			}
		}

		return trueskill.Calculate(retval), nil
	} else {
		return rs, nil
	}
}

func CreatePlay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs models.Play

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		err := fn(args.Data, true)
		if err != nil {
			return nil, err
		}
	}

	query, err := args.Insert("tboardgameplays")
	if err != nil {
		return nil, err
	}

	row, err := db.NamedExec(query.String(), &rs)
	if err != nil {
		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	rs.Id = id
	return rs, nil
}

func UpdatePlay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getPlay(db, args.Id)
	if err != nil {
		return nil, err
	}

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		fn(args.Data, false)
	}

	sql, err := args.Update("tboardgameplays")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func DeletePlay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getPlay(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgameplays where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
