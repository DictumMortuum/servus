package boardgames

import (
	"database/sql"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func createPrice(db *sqlx.DB, data models.BoardgamePrice) (int64, error) {
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

func updatePrice(db *sqlx.DB, data models.BoardgamePrice) error {
	sql := `
	update tboardgameprices set
		cr_date = :cr_date,
		date = NOW(),
		boardgame = :boardgame,
		store = :store,
		original_price = :original_price,
		reduced_price = :reduced_price,
		price_diff = :price_diff,
		seq = seq + 1
	where id = :id`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func priceExists(db *sqlx.DB, row models.BoardgamePrice) (int64, error) {
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
		reduced_price = :reduced_price
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

func sendTextForPrice(db *sqlx.DB, data models.BoardgamePrice) error {
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

func getPricesWithoutTexts(db *sqlx.DB) ([]models.BoardgamePrice, error) {
	sql := `
	select
		*
	from
		tboardgameprices
	where
		text_send = 0`

	return getPrices(db, sql)
}

func getPrices(db *sqlx.DB, sql string) ([]models.BoardgamePrice, error) {
	rs := []models.BoardgamePrice{}

	err := db.Select(&rs, sql)
	if err != nil {
		return rs, err
	}

	return rs, nil
}
