package models

import (
	"fmt"
	"time"
)

type BoardgamePrice struct {
	Id            int64         `db:"id"`
	CrDate        time.Time     `db:"cr_date"`
	BoardgameId   JsonNullInt64 `db:"boardgame_id"`
	Boardgame     string        `db:"boardgame" json:"boardgame"`
	StoreId       JsonNullInt64 `db:"store_id"`
	Store         string        `db:"store" json:"store"`
	OriginalPrice float64       `db:"original_price"`
	ReducedPrice  float64       `db:"reduced_price"`
}

func (obj BoardgamePrice) Msg() string {
	return fmt.Sprintf("%s offers %s at %.2f from %.2f\n", obj.Store, obj.Boardgame, obj.ReducedPrice, obj.OriginalPrice)
}

func (obj BoardgamePrice) Insert() string {
	return `
	insert into tboardgameprices (
		cr_date,
		boardgame_id,
		store_id,
		original_price,
		reduced_price
	) values (
		NOW(),
		:boardgame_id,
		:store_id,
		:original_price,
		:reduced_price
	)`
}

func (obj BoardgamePrice) Exists() string {
	return "select id from tboardgameprices where boardgame_id = :boardgame_id and store_id = :store_id"
}
