package models

type Insertable interface {
	Insert() string
	Exists() string
}

type Selectable interface {
	Select() string
}

type Sendable interface {
	Msg() string
}
