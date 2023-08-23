package models

import (
	"errors"
	"time"
)

type Boardgame struct {
	Id             int64           `db:"id" json:"id"`
	Name           string          `db:"name" json:"name"`
	Date           time.Time       `json:"validUntil"`
	UpdatedAt      time.Time       `db:"updated_at"`
	Data           Json            `db:"data" json:"data"`
	BggData        Json            `db:"bgg_data" json:"bgg_data"`
	Guid           JsonNullString  `db:"tx_guid" json:"tx_guid"`
	Cost           JsonNullFloat64 `db:"cost" json:"cost"`
	Rank           JsonNullInt64   `db:"rank" json:"rank"`
	Thumb          JsonNullString  `db:"thumb" json:"thumb"`
	Preview        JsonNullString  `db:"preview" json:"preview"`
	Configured     bool            `db:"configured" json:"configured"`
	BggDataNotNull bool            `db:"bgg_data_not_null" json:"bgg_data_not_null"`
	RankNotNull    bool            `db:"rank_not_null" json:"rank_not_null"`
	Year           JsonNullInt64   `db:"year" json:"year"`
	MinPlayers     JsonNullInt64   `db:"min_players" json:"minplayers"`
	MaxPlayers     JsonNullInt64   `db:"max_players" json:"maxplayers"`
	Square200      JsonNullString  `db:"square200" json:"square200"`
}

func (rs *Boardgame) SetName(data map[string]interface{}, create bool) error {
	name, err := getString(data, "name")
	if err != nil {
		return err
	}

	rs.Name = name
	return nil
}

func (rs *Boardgame) SetData(data map[string]interface{}, create bool) error {
	if val, ok := data["data"]; ok {
		return rs.Data.Scan(val)
	}

	if create {
		rs.Data = nil
		return nil
	}

	return errors.New("could not find 'data' parameter")
}

func (rs *Boardgame) SetRank(data map[string]interface{}, create bool) error {
	rank, err := getJsonNullInt64(data, "rank")
	if create || err == nil {
		rs.Rank = rank
		return nil
	}

	return err
}

func (rs *Boardgame) SetThumb(data map[string]interface{}, create bool) error {
	thumb, err := getJsonNullString(data, "thumb")
	if create || err == nil {
		rs.Thumb = thumb
		return nil
	}

	return err
}

func (rs *Boardgame) SetPreview(data map[string]interface{}, create bool) error {
	preview, err := getJsonNullString(data, "preview")
	if create || err == nil {
		rs.Preview = preview
		return nil
	}

	return err
}

func (rs *Boardgame) SetConfigured(data map[string]interface{}, create bool) error {
	configured, err := getBool(data, "configured")
	if create || err == nil {
		rs.Configured = configured
		return nil
	}

	return err
}

func (rs *Boardgame) SetId(data map[string]interface{}, create bool) error {
	if create {
		id, err := getInt64(data, "id")
		if err != nil {
			return err
		}

		rs.Id = id
		return nil
	}

	return nil
}

func (rs *Boardgame) Constructor() []func(map[string]interface{}, bool) error {
	return []func(map[string]interface{}, bool) error{
		rs.SetId,
		rs.SetName,
		rs.SetData,
		rs.SetRank,
		rs.SetThumb,
		rs.SetPreview,
		rs.SetConfigured,
	}
}
