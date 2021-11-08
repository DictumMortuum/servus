package models

import (
	"time"
)

type Boardgame struct {
	Id   int64     `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
	Date time.Time `json:"validUntil"`
}

func (obj Boardgame) Insert() string {
	return `insert into tboardgames (id,name) values (:id,:name)`
}

func (obj Boardgame) Select() string {
	return `select * from tboardgames where name = :name`
}

func (obj Boardgame) Exists() string {
	return `select id from tboardgames where name = :name`
}
