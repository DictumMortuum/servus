package boardgames

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	trueskill "github.com/mafredri/go-trueskill"
	"sort"
	"text/template"
	"time"
)

type Play struct{}

func getPlay(db *sqlx.DB, id int64) (*models.Play, error) {
	var rs models.Play

	err := db.QueryRowx(`
		select
			p.*,
			g.name
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
	scoreFunc, sortFunc := getFuncs(play.BoardgameId)
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

func calculateTrueskill(plays []models.Play) []models.Play {
	ts := trueskill.New()
	players := map[string]trueskill.Player{}

	for idx := range plays {
		// Reverse the array, so that winner is on the top.
		for i, j := 0, len(plays[idx].Stats)-1; i < j; i, j = i+1, j-1 {
			plays[idx].Stats[i], plays[idx].Stats[j] = plays[idx].Stats[j], plays[idx].Stats[i]
		}

		// Generate a list of the names of the players
		playersInPlay := []string{}
		for _, stat := range plays[idx].Stats {
			playersInPlay = append(playersInPlay, stat.Player)
		}

		// Generate a list of draws, if there are any
		draws := []bool{}
		flag := false
		for i := 0; i < len(plays[idx].Stats)-1; i++ {
			if plays[idx].Stats[i].Data["score"] == plays[idx].Stats[i+1].Data["score"] {
				draws = append(draws, true)
				flag = true
			} else {
				draws = append(draws, false)
			}
		}

		if flag {
			// fmt.Println(plays[idx])
			// fmt.Println(draws)
		}

		// For each player that participated, generate a trueskill structure - initialize it if it's his first time.
		var playerSkills []trueskill.Player
		for _, name := range playersInPlay {
			if val, ok := players[name]; ok {
				playerSkills = append(playerSkills, trueskill.NewPlayer(val.Mu(), val.Sigma()))
			} else {
				players[name] = ts.NewPlayer()
				playerSkills = append(playerSkills, players[name])
			}
		}

		adjustedPlayers, probability := ts.AdjustSkillsWithDraws(playerSkills, draws)
		// adjustedPlayers and player names are in order - useful to copy over the new ratings
		for i, name := range playersInPlay {
			players[name] = trueskill.NewPlayer(adjustedPlayers[i].Mu(), adjustedPlayers[i].Sigma())
			plays[idx].Stats[i].Mu = adjustedPlayers[i].Mu()
			plays[idx].Stats[i].Sigma = adjustedPlayers[i].Sigma()
			plays[idx].Stats[i].TrueSkill = ts.TrueSkill(players[name])
		}
		plays[idx].Probability = probability * 100
	}

	return plays
}

func (obj Play) Get(db *sqlx.DB, id int64) (interface{}, error) {
	return getPlay(db, id)
}

func (obj Play) GetList(db *sqlx.DB, args models.QueryBuilder) (interface{}, int, error) {
	var rs []models.Play

	var count []int
	err := db.Select(&count, "select 1 from tboardgameplays")
	if err != nil {
		return nil, -1, err
	}

	sql := `
		select
			p.*,
			g.name
		from
			tboardgameplays p,
			tboardgames g
		where
			p.boardgame_id = g.id
		{{ if gt (len .Id) 0 }}
			and p.{{ .RefKey }} in (?)
		{{ else if gt (len .FilterVal) 0 }}
			and p.{{ .FilterKey }} = "{{ .FilterVal }}"
		{{ end }}
		{{ if gt (len .Sort) 0 }}
		order by p.{{ .Sort }} {{ .Order }}
		{{ else }}
		order by p.date asc
		{{ end }}
		{{ if eq (len .Range) 2 }}
		limit {{ index .Range 0 }}, {{ .Page }}
		{{ end }}`

	var tpl bytes.Buffer
	t := template.Must(template.New("list").Parse(sql))
	err = t.Execute(&tpl, args)
	if err != nil {
		return nil, -1, err
	}

	query, ids, err := sqlx.In(tpl.String(), args.Id)
	if err != nil {
		query = tpl.String()
	} else {
		query = db.Rebind(query)
	}

	err = db.Select(&rs, query, ids...)
	if err != nil {
		return nil, -1, err
	}

	if args.Resources["stats"] {
		retval := []models.Play{}

		for _, item := range rs {
			stats, err := getPlayStats(db, item.Id)
			if err != nil {
				return nil, -1, err
			}
			item.Stats = stats

			play, err := scorePlay(item)
			if err != nil {
				return nil, -1, err
			}

			retval = append(retval, *play)
		}

		return calculateTrueskill(retval), len(count), nil
	} else {
		return rs, len(count), nil
	}
}

func (obj Play) Create(db *sqlx.DB, qb models.QueryBuilder) (interface{}, error) {
	var rs models.Play

	if val, ok := qb.Data["date"]; ok {
		//"Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"
		t, err := time.Parse("2006-01-02", val.(string))
		if err != nil {
			return nil, err
		}

		rs.Date = t
	} else {
		rs.Date = time.Now()
	}

	if val, ok := qb.Data["boardgame_id"]; ok {
		rs.BoardgameId = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'boardgame_id' parameter")
	}

	rs.CrDate = time.Now()

	query, err := qb.Insert("tboardgameplays")
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

func (obj Play) Update(db *sqlx.DB, id int64, qb models.QueryBuilder) (interface{}, error) {
	rs, err := getPlay(db, id)
	if err != nil {
		return nil, err
	}

	sql, err := qb.Update("tboardgameplays")
	if err != nil {
		return nil, err
	}

	if val, ok := qb.Data["date"]; ok {
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

func (obj Play) Delete(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getPlay(db, id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgameplays where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
