package models

type Store struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (obj Store) Insert() string {
	return `insert into tboardgamestores (name) values (:name)`
}

func (obj Store) Exists() string {
	return `select id from tboardgamestores where name = :name`
}
