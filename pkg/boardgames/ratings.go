package boardgames

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus/pkg/boardgames/score"
	"github.com/DictumMortuum/servus/pkg/boardgames/trueskill"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func getBoardgamesWithPlays(db *sqlx.DB) ([]int64, error) {
	var rs []int64

	sql := `
		select
			distinct boardgame_id
		from
			tboardgameplays
	`

	err := db.Select(&rs, sql)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getPlays(db *sqlx.DB) ([]models.Play, error) {
	var rs []models.Play

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
		order by p.date asc, p.id
	`

	err := db.Select(&rs, sql)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getPlaysForBoardgame(db *sqlx.DB, id int64) ([]models.Play, error) {
	var rs []models.Play

	sql := `
		select
			p.*,
			g.name,
			g.data
		from
			tboardgameplays p,
			tboardgames g
		where
			p.boardgame_id = ? and
			p.boardgame_id = g.id
		order by p.date asc, p.id
	`

	err := db.Select(&rs, sql, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func playsAreCooperative(plays []models.Play) bool {
	count := 0

	for _, play := range plays {
		if play.IsCooperative() {
			count += 1
		}
	}

	return len(plays) == count
}

type Stats struct {
	PlayerId      int64   `json:"player_id"`
	Player        string  `json:"player"`
	PlayerSurname string  `json:"player_surname"`
	Mu            float64 `json:"mu"`
	Sigma         float64 `json:"sigma"`
	TrueSkill     float64 `json:"trueskill"`
}

type list struct {
	Id    int64                 `json:"id"`
	Name  string                `json:"name"`
	Thumb models.JsonNullString `json:"thumb"`
	Count int                   `json:"count"`
	List  []Stats               `json:"ratings"`
}

func getLatestScore(plays []models.Play) []Stats {
	rs := map[int64]Stats{}

	for _, play := range plays {
		for _, stat := range play.Stats {
			id := stat.PlayerId
			rs[id] = Stats{
				stat.PlayerId,
				stat.Player,
				stat.PlayerSurname,
				stat.Mu,
				stat.Sigma,
				stat.TrueSkill,
			}
		}
	}

	tmp := []Stats{}

	for _, val := range rs {
		tmp = append(tmp, val)
	}

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].TrueSkill > tmp[j].TrueSkill
	})

	return tmp
}

func GetTrueskillLists(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []list{}

	boardgames, err := getBoardgamesWithPlays(db)
	if err != nil {
		return nil, err
	}

	for _, id := range boardgames {
		boardgame, err := getBoardgame(db, id)
		if err != nil {
			return nil, err
		}

		plays, err := getPlaysForBoardgame(db, id)
		if err != nil {
			return nil, err
		}

		if playsAreCooperative(plays) {
			continue
		}

		if len(plays) < 5 {
			continue
		}

		scored_plays, err := score.CalculateAll(db, plays)
		if err != nil {
			return nil, err
		}

		trueskill_plays := trueskill.Calculate(scored_plays)

		rs = append(rs, list{
			Id:    boardgame.Id,
			Name:  boardgame.Name,
			Thumb: boardgame.Thumb,
			Count: len(scored_plays),
			List:  getLatestScore(trueskill_plays),
		})
	}

	return rs, nil
}

func GetTrueskillOverall(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []list{}

	plays, err := getPlays(db)
	if err != nil {
		return nil, err
	}

	scored_plays, err := score.CalculateAll(db, plays)
	if err != nil {
		return nil, err
	}

	trueskill_plays := trueskill.Calculate(scored_plays)

	for _, play := range trueskill_plays {
		players := []int64{}
		for _, stat := range play.Stats {
			players = append(players, stat.PlayerId)
		}

		var teams [][]int64
		if val, ok := play.PlaySettings["teams2"]; ok {
			teams = val.([][]int64)
		}

		var winners []int64
		if len(teams) > 0 {
			winners = teams[0]
		} else {
			winners = []int64{play.Stats[0].PlayerId}
			for i, draw := range play.Draws {
				if draw {
					winners = append(winners, play.Stats[i+1].PlayerId)
				} else {
					break
				}
			}
		}

		payload := map[string]interface{}{
			"winners": winners,
			"players": players,
			"draws":   play.Draws,
		}

		if len(teams) > 0 {
			payload["teams"] = teams
		}

		json, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		sql := fmt.Sprintf(`
			update tboardgameplays set play_data = '%s' where id = ?
		`, string(json))

		_, err = db.Exec(sql, play.Id)
		if err != nil {
			return nil, err
		}
	}

	for _, play := range scored_plays {
		if !play.IsCooperative() {
			continue
		}

		players := []int64{}
		for _, stat := range play.Stats {
			players = append(players, stat.PlayerId)
		}

		winners := []int64{}
		for i, stat := range play.Stats {
			if val, ok := stat.Data["won"]; ok {
				if val.(bool) {
					winners = append(winners, play.Stats[i].PlayerId)
				}
			}
		}

		payload := map[string]interface{}{
			"cooperative": true,
			"players":     players,
			"winners":     winners,
		}

		json, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		sql := fmt.Sprintf(`
			update tboardgameplays set play_data = '%s' where id = ?
		`, string(json))

		_, err = db.Exec(sql, play.Id)
		if err != nil {
			return nil, err
		}
	}

	rs = append(rs, list{
		Name:  "Overall",
		Count: len(scored_plays),
		List:  getLatestScore(trueskill_plays),
	})

	return rs, nil
}
