package models

type Book struct {
	Id        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Isbn      string `db:"isbn" json:"isbn"`
	Publisher string `db:"publisher" json:"publisher"`
	Writer    string `db:"writer" json:"writer"`
}

func (rs *Book) SetName(data map[string]interface{}, create bool) error {
	name, err := getString(data, "name")
	if err != nil {
		return err
	}

	rs.Name = name
	return nil
}

func (rs *Book) SetWriter(data map[string]interface{}, create bool) error {
	writer, err := getString(data, "writer")
	if err != nil {
		return err
	}

	rs.Writer = writer
	return nil
}

func (rs *Book) SetIsbn(data map[string]interface{}, create bool) error {
	isbn, err := getString(data, "isbn")
	if err != nil {
		return err
	}

	rs.Isbn = isbn
	return nil
}

func (rs *Book) SetPublisher(data map[string]interface{}, create bool) error {
	publisher, err := getString(data, "publisher")
	if err != nil {
		return err
	}

	rs.Publisher = publisher
	return nil
}

func (rs *Book) Constructor() []func(map[string]interface{}, bool) error {
	return []func(map[string]interface{}, bool) error{
		rs.SetName,
		rs.SetIsbn,
		rs.SetPublisher,
		rs.SetWriter,
	}
}
