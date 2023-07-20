package models

import (
	"errors"
)

type Stats struct {
	Id            int64   `db:"id" json:"id"`
	PlayId        int64   `db:"play_id" json:"play_id"`
	BoardgameId   int64   `db:"boardgame_id" json:"boardgame_id"`
	PlayerId      int64   `db:"player_id" json:"player_id"`
	Data          Json    `db:"data" json:"data"`
	Player        string  `db:"name" json:"player"`
	PlayerSurname string  `db:"surname" json:"player_surname"`
	Mu            float64 `json:"mu"`
	Sigma         float64 `json:"sigma"`
	TrueSkill     float64 `json:"trueskill"`
	Delta         float64 `json:"delta"`
}

func (rs *Stats) Constructor() []func(map[string]interface{}, bool) error {
	return []func(map[string]interface{}, bool) error{
		rs.SetPlayId,
		rs.SetBoardgameId,
		rs.SetPlayerId,
		rs.SetData,
	}
}

func (rs *Stats) SetPlayId(data map[string]interface{}, create bool) error {
	id, err := getInt64(data, "play_id")
	if err != nil {
		return err
	}

	rs.PlayId = id
	return nil
}

func (rs *Stats) SetBoardgameId(data map[string]interface{}, create bool) error {
	id, err := getInt64(data, "boardgame_id")
	if err != nil {
		return err
	}

	rs.BoardgameId = id
	return nil
}

func (rs *Stats) SetPlayerId(data map[string]interface{}, create bool) error {
	id, err := getInt64(data, "player_id")
	if err != nil {
		return err
	}

	rs.PlayerId = id
	return nil
}

func (rs *Stats) SetData(data map[string]interface{}, create bool) error {
	if val, ok := data["data"]; ok {
		return rs.Data.Scan(val)
	}

	if create {
		rs.Data = nil
		return nil
	}

	return errors.New("could not find 'data' parameter")
}
