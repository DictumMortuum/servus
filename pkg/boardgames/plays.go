package boardgames

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/boardgames/trueskill"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"sort"
	"text/template"
	"time"
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

	play, err := scorePlay(rs)
	if err != nil {
		return nil, err
	}

	return play, nil
}

func scorePlay(play models.Play) (*models.Play, error) {
	if play.IsCooperative() {
		return &play, nil
	}

	scoreFunc, sortFunc := getFuncs(play)
	if scoreFunc == nil || sortFunc == nil {
		e := fmt.Sprintf("Could not find sort or score function for boardgame %s\n", play.Boardgame)
		return nil, errors.New(e)
	}

	rs := []models.Stats{}
	for _, item := range play.Stats {
		item.Data["score"] = scoreFunc(item)
		rs = append(rs, item)
	}
	sort.Slice(rs, sortFunc(rs))

	play.Stats = rs

	return &play, nil
}

func GetPlay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getPlay(db, args.Id)
}

func GetListPlay(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Play

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
			stats, err := getPlayStats(db, item.Id)
			if err != nil {
				return nil, err
			}
			item.Stats = stats

			play, err := scorePlay(item)
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

	if val, ok := args.Data["date"]; ok {
		//"Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"
		t, err := time.Parse("2006-01-02", val.(string))
		if err != nil {
			return nil, err
		}

		rs.Date = t
	} else {
		rs.Date = time.Now()
	}

	if val, ok := args.Data["boardgame_id"]; ok {
		rs.BoardgameId = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'boardgame_id' parameter")
	}

	rs.CrDate = time.Now()

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

	sql, err := args.Update("tboardgameplays")
	if err != nil {
		return nil, err
	}

	if val, ok := args.Data["date"]; ok {
		t, err := time.Parse("2006-01-02T15:04:05-0700", val.(string))
		if err != nil {
			return nil, err
		}

		rs.Date = t
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
