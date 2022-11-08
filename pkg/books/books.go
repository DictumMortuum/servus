package books

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func getBook(db *sqlx.DB, id int64) (*models.Book, error) {
	var rs models.Book

	err := db.QueryRowx(`select * from tbooks where id = ?`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func GetBook(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getBook(db, args.Id)
}

func GetListBook(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []models.Book{}

	var count []int
	err := db.Select(&count, "select 1 from tbooks")
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", len(count)))

	sql, err := args.List(`
		select * from tbooks
	`)
	if err != nil {
		return nil, err
	}

	query, ids, err := sqlx.In(sql.String(), args.Ids)
	if err != nil {
		query = sql.String()
	} else {
		query = db.Rebind(query)
	}

	err = db.Select(&rs, query, ids...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func CreateBook(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs models.Book

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		err := fn(args.Data, true)
		if err != nil {
			return nil, err
		}
	}

	query, err := args.Insert("tboardgameBooks")
	if err != nil {
		return nil, err
	}

	row, err := db.NamedExec(query.String(), &rs)
	if err != nil {
		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	rs.Id = id
	return rs, nil
}

func UpdateBook(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getBook(db, args.Id)
	if err != nil {
		return nil, err
	}

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		fn(args.Data, false)
	}

	sql, err := args.Update("tbooks")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func DeleteBook(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getBook(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tbooks where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
