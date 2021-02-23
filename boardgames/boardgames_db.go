package boardgames

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type PlayersRow struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type PriceRow struct {
	Id            int64     `db:"id"`
	CrDate        time.Time `db:"cr_date"`
	Date          time.Time `db:"date"`
	Boardgame     string    `db:"boardgame"`
	Store         string    `db:"store"`
	OriginalPrice float64   `db:"original_price"`
	ReducedPrice  float64   `db:"reduced_price"`
	PriceDiff     float64   `db:"price_diff"`
	Link          string    `db:"link"`
	TextSend      bool      `db:"text_send"`
	Seq           int       `db:"seq"`
}

func CreatePlayer(db *sqlx.DB, data PlayersRow) error {
	sql := `
	insert into tboardgameplayers (
		name
	) values (
		:name
	)`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func GetPlayer(db *sqlx.DB, id int64) (*PlayersRow, error) {
	var retval PlayersRow

	err := db.QueryRowx(`select * from tboardgameplayers where id = ?`, id).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}

func GetPlayerByName(db *sqlx.DB, name string) (*PlayersRow, error) {
	var retval PlayersRow

	err := db.QueryRowx(`select * from tboardgameplayers where name = ?`, name).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}

func CreatePrice(db *sqlx.DB, data PriceRow) (int64, error) {
	sql := `
	insert into tboardgameprices (
		cr_date,
		date,
		boardgame,
		store,
		original_price,
		reduced_price,
		price_diff,
		link,
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
		NULL,
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

func UpdatePrice(db *sqlx.DB, data PriceRow) error {
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

func PriceExists(db *sqlx.DB, row PriceRow) (int64, error) {
	var id sql.NullInt64

	sql := `
	select
		id
	from
		tboardgameprices
	where
		boardgame = :boardgame and
		store = :store and
		original_price - :original_price >= -1 and original_price - :original_price < 1
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
