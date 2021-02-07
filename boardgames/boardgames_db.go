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

func CreatePrice(db *sqlx.DB, data PriceRow) error {
	sql := `
	insert into tboardgameprices (
		cr_date,
		date,
		boardgame,
		store,
		original_price,
		reduced_price,
		price_diff,
		link
	) values (
		NOW(),
		:date,
		:boardgame,
		:store,
		:original_price,
		:reduced_price,
		:price_diff,
		NULL
	)`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
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
		price_diff = :price_diff
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
		boardgame = :boardgame
	and store = :store
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
