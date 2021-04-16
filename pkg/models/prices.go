package models

import (
	"fmt"
	"time"
)

type BoardgamePrice struct {
	Id            int64     `db:"id"`
	CrDate        time.Time `db:"cr_date"`
	Date          time.Time `db:"date"`
	Boardgame     string    `db:"boardgame"`
	Store         string    `db:"store"`
	OriginalPrice float64   `db:"original_price"`
	ReducedPrice  float64   `db:"reduced_price"`
	PriceDiff     float64   `db:"price_diff"`
	TextSend      bool      `db:"text_send"`
	Seq           int       `db:"seq"`
}

func (obj BoardgamePrice) Msg() string {
	return fmt.Sprintf("%s offers %s at %.2f from %.2f\n", obj.Store, obj.Boardgame, obj.ReducedPrice, obj.OriginalPrice)
}
