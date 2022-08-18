package models

type Store struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (rs *Store) Constructor() []func(map[string]interface{}, bool) error {
	return []func(map[string]interface{}, bool) error{
		rs.SetName,
	}
}

func (rs *Store) SetName(data map[string]interface{}, create bool) error {
	name, err := getString(data, "name")
	if err != nil {
		return err
	}

	rs.Name = name
	return nil
}
