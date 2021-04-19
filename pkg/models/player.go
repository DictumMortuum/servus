package models

type Player struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}
