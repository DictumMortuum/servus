package models

type Player struct {
	Id      int64          `db:"id" json:"id"`
	Name    string         `db:"name" json:"name"`
	Surname string         `db:"surname" json:"surname"`
	Email   JsonNullString `db:"email" json:"email"`
	Uuid    JsonNullString `db:"uuid" json:"uuid"`
}

func (rs *Player) SetName(data map[string]interface{}, create bool) error {
	name, err := getString(data, "name")
	if err != nil {
		return err
	}

	rs.Name = name
	return nil
}

func (rs *Player) SetSurname(data map[string]interface{}, create bool) error {
	surname, err := getString(data, "surname")
	if err != nil {
		return err
	}

	rs.Surname = surname
	return nil
}

func (rs *Player) SetEmail(data map[string]interface{}, create bool) error {
	email, err := getString(data, "email")
	if err != nil {
		return err
	}

	rs.Email = NewJsonNullString(email)
	return nil
}

func (rs *Player) Constructor() []func(map[string]interface{}, bool) error {
	return []func(map[string]interface{}, bool) error{
		rs.SetName,
		rs.SetSurname,
		rs.SetEmail,
	}
}
