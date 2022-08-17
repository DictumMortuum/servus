package models

type Mapping struct {
	Id          int64  `db:"id" json:"id"`
	BoardgameId int64  `db:"boardgame_id" json:"boardgame_id"`
	Name        string `db:"name" json:"name"`
}

func (rs *Mapping) SetName(data map[string]interface{}, create bool) error {
	name, err := getString(data, "name")
	if err != nil {
		return err
	}

	rs.Name = name
	return nil
}

func (rs *Mapping) SetBoardgameId(data map[string]interface{}, create bool) error {
	id, err := getInt64(data, "boardgame_id")
	if err != nil {
		return err
	}

	rs.BoardgameId = id
	return nil
}

func (rs *Mapping) Constructor() []func(map[string]interface{}, bool) error {
	return []func(map[string]interface{}, bool) error{
		rs.SetBoardgameId,
		rs.SetName,
	}
}
