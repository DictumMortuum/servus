package models

import (
	"github.com/jmoiron/sqlx"
)

type Getable interface {
	Get(*sqlx.DB, int64) (interface{}, error)
}

type Getlistable interface {
	GetTable() string
	GetList(*sqlx.DB, string, ...interface{}) (interface{}, error)
}

type Createable interface {
	GetTable() string
	Create(*sqlx.DB, string, map[string]interface{}) (interface{}, error)
}

type Updateable interface {
	GetTable() string
	Update(*sqlx.DB, string, int64, map[string]interface{}) (interface{}, error)
}

type Deleteable interface {
	GetTable() string
	Delete(*sqlx.DB, string, int64) (interface{}, error)
}
