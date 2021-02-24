package boardgames

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type PriceRow struct {
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

func (p PriceRow) Msg() string {
	return fmt.Sprintf("%s offers %s at %.2f from %.2f\n", p.Store, p.Boardgame, p.ReducedPrice, p.OriginalPrice)
}

func createPrice(db *sqlx.DB, data PriceRow) (int64, error) {
	sql := `
	insert into tboardgameprices (
		cr_date,
		date,
		boardgame,
		store,
		original_price,
		reduced_price,
		price_diff,
		text_send,
		seq
	) values (
		NOW(),
		:date,
		:boardgame,
		:store,
		:original_price,
		:reduced_price,
		:price_diff,
		0,
		0
	)`

	res, err := db.NamedExec(sql, &data)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func updatePrice(db *sqlx.DB, data PriceRow) error {
	sql := `
	update tboardgameprices set
		cr_date = NOW(),
		date = :date,
		boardgame = :boardgame,
		store = :store,
		original_price = :original_price,
		reduced_price = :reduced_price,
		price_diff = :price_diff,
		text_send = 0,
		seq = seq + 1
	where id = :id`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func priceExists(db *sqlx.DB, row PriceRow) (int64, error) {
	var id sql.NullInt64

	sql := `
	select
		id
	from
		tboardgameprices
	where
		boardgame = :boardgame and
		store = :store and
		original_price = :original_price and
		reduced_price - :reduce_price
	`

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return -1, err
	}

	err = stmt.Get(&id, row)
	if err != nil {
		return -1, nil
	}

	retval, err := id.Value()
	if err != nil {
		return -1, err
	}

	if retval == nil {
		return -1, err
	}

	return retval.(int64), nil
}

func sendTextForPrice(db *sqlx.DB, data PriceRow) error {
	sql := `
	update
		tboardgameprices
	set
		text_send = 1
	where
		id = :id`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func getPricesWithoutTexts(db *sqlx.DB) ([]PriceRow, error) {
	sql := `
	select
		*
	from
		tboardgameprices
	where
		text_send = 0`

	return getPrices(db, sql)
}

func getPrices(db *sqlx.DB, sql string) ([]PriceRow, error) {
	rs := []PriceRow{}

	err := db.Select(&rs, sql)
	if err != nil {
		return rs, err
	}

	return rs, nil
}
