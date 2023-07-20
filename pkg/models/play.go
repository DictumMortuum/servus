package models

import (
	"time"
)

type Play struct {
	Id                int64         `db:"id" json:"id"`
	BoardgameId       int64         `db:"boardgame_id" json:"boardgame_id"`
	LocationId        JsonNullInt64 `db:"location_id" json:"location_id"`
	BoardgameSettings Json          `db:"data" json:"boardgame_data"`
	PlaySettings      Json          `db:"play_data" json:"play_data"`
	CrDate            time.Time     `db:"cr_date" json:"cr_date"`
	Date              time.Time     `db:"date" json:"date"`
	Boardgame         string        `db:"name" json:"boardgame"`
	Stats             []Stats       `json:"stats"`
	Probability       float64       `json:"probability"`
	Draws             []bool        `json:"draws"`
	// PlayData          Json          `db:"play_data" json:"play_data"`
}

func (p Play) IsCooperative() bool {
	if val, ok := p.BoardgameSettings["cooperative"]; ok {
		return val.(bool)
	}

	return false
}

func (p Play) GetTeams() [][]int64 {
	rs := [][]int64{}

	if val, ok := p.PlaySettings["teams"]; ok {
		raw_teams := val.([]interface{})

		for _, raw_team := range raw_teams {
			team := []int64{}

			for _, item := range raw_team.([]interface{}) {
				team = append(team, int64(item.(float64)))
			}

			rs = append(rs, team)
		}
	}

	return rs
}

func (rs *Play) SetDate(data map[string]interface{}, create bool) error {
	date, err := getTime(data, "date")
	if err != nil {
		return err
	}

	rs.Date = date
	return nil
}

func (rs *Play) SetCrDate(data map[string]interface{}, create bool) error {
	rs.CrDate = time.Now()
	return nil
}

func (rs *Play) SetBoardgameId(data map[string]interface{}, create bool) error {
	id, err := getInt64(data, "boardgame_id")
	if err != nil {
		return err
	}

	rs.BoardgameId = id
	return nil
}

func (rs *Play) Constructor() []func(map[string]interface{}, bool) error {
	return []func(map[string]interface{}, bool) error{
		rs.SetBoardgameId,
		rs.SetDate,
		rs.SetCrDate,
	}
}
