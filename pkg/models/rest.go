package models

import (
	"github.com/jmoiron/sqlx"
)

type Getable interface {
	Get(*sqlx.DB, int64) (interface{}, error)
}

type Getlistable interface {
	GetList(*sqlx.DB, QueryBuilder) (interface{}, int, error)
}

type Createable interface {
	Create(*sqlx.DB, QueryBuilder) (interface{}, error)
}

type Updateable interface {
	Update(*sqlx.DB, int64, QueryBuilder) (interface{}, error)
}

type Deleteable interface {
	Delete(*sqlx.DB, int64) (interface{}, error)
}
