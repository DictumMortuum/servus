package boardgames

import (
	"bytes"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"text/template"
)

func GetHistoricPriceById(db *sqlx.DB, id int64) (*models.HistoricPrice, error) {
	var rs models.HistoricPrice

	err := db.QueryRowx(`
		select
			p.*
		from
			tboardgamepriceshistory p
		where
			p.id = ?
	`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func GetHistoricPrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return GetHistoricPriceById(db, args.Id)
}

func countHistoricPrices(db *sqlx.DB, args *models.QueryBuilder) (int, error) {
	var tpl bytes.Buffer

	sql := `
		select
			1
		from
			tboardgamepriceshistory p
	`

	t := template.Must(template.New("count").Parse(sql))
	err := t.Execute(&tpl, args)
	if err != nil {
		return -1, err
	}

	var count []int
	err = db.Select(&count, tpl.String())
	if err != nil {
		return -1, err
	}

	return len(count), nil
}

func GetListHistoricPrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.HistoricPrice

	count, err := countHistoricPrices(db, args)
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", count))

	sql := `
		select
			*
		from
			tboardgamepriceshistory p
	`

	var tpl bytes.Buffer
	t := template.Must(template.New("list").Parse(sql))
	err = t.Execute(&tpl, args)
	if err != nil {
		return nil, err
	}

	query, ids, err := sqlx.In(tpl.String(), args.Ids)
	if err != nil {
		query = tpl.String()
	} else {
		query = db.Rebind(query)
	}

	err = db.Select(&rs, query, ids...)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

// func CreatePrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
// 	var rs models.Price

// 	if val, ok := args.Data["boardgame_id"]; ok {
// 		rs.BoardgameId = models.JsonNullInt64{
// 			Int64: int64(val.(float64)),
// 			Valid: true,
// 		}
// 	} else {
// 		return nil, errors.New("please provide a 'boardgame_id' parameter")
// 	}

// 	if val, ok := args.Data["store_id"]; ok {
// 		rs.StoreId = int64(val.(float64))
// 	} else {
// 		return nil, errors.New("please provide a 'store_id' parameter")
// 	}

// 	if val, ok := args.Data["price"]; ok {
// 		rs.Price = val.(float64)
// 	} else {
// 		return nil, errors.New("please provide a 'price' parameter")
// 	}

// 	if val, ok := args.Data["stock"]; ok {
// 		rs.Stock = val.(bool)
// 	} else {
// 		return nil, errors.New("please provide a 'stock' parameter")
// 	}

// 	if val, ok := args.Data["url"]; ok {
// 		rs.Url = val.(string)
// 	} else {
// 		return nil, errors.New("please provide a 'url' parameter")
// 	}

// 	if val, ok := args.Data["store_thumb"]; ok {
// 		rs.StoreThumb = val.(string)
// 	} else {
// 		return nil, errors.New("please provide a 'store_thumb' parameter")
// 	}

// 	if val, ok := args.Data["batch"]; ok {
// 		rs.Batch = val.(int64)
// 	} else {
// 		return nil, errors.New("please provide a 'batch' parameter")
// 	}

// 	rs.Levenshtein = 0
// 	rs.Hamming = 0

// 	query, err := args.Insert("tboardgameprices")
// 	if err != nil {
// 		return nil, err
// 	}

// 	price, err := db.NamedExec(query.String(), &rs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	id, err := price.LastInsertId()
// 	if err != nil {
// 		return nil, err
// 	}

// 	rs.Id = id
// 	return rs, nil
// }

// func UpdatePrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
// 	rs, err := GetPriceById(db, args.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	args.IgnoreColumn("rank")
// 	args.IgnoreColumn("thumb")
// 	args.IgnoreColumn("boardgame_name")
// 	args.IgnoreColumn("store_name")

// 	if val, ok := args.Data["boardgame_id"]; ok {
// 		rs.BoardgameId = models.JsonNullInt64{
// 			Int64: int64(val.(float64)),
// 			Valid: true,
// 		}

// 		exists, err := boardgameExists(db, map[string]interface{}{
// 			"id": rs.BoardgameId.Int64,
// 		})
// 		if err != nil {
// 			return nil, err
// 		}

// 		if exists == nil {
// 			_, err := bgg.FetchBoardgame(db, rs.BoardgameId.Int64)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}

// 	if val, ok := args.Data["store_id"]; ok {
// 		rs.StoreId = int64(val.(float64))
// 	}

// 	if val, ok := args.Data["store_thumb"]; ok {
// 		rs.StoreThumb = val.(string)
// 	}

// 	if val, ok := args.Data["price"]; ok {
// 		rs.Price = val.(float64)
// 	}

// 	if val, ok := args.Data["stock"]; ok {
// 		rs.Stock = val.(bool)
// 	}

// 	if val, ok := args.Data["url"]; ok {
// 		rs.Url = val.(string)
// 	}

// 	if val, ok := args.Data["batch"]; ok {
// 		rs.Batch = int64(val.(float64))
// 	}

// 	if val, ok := args.Data["mapped"]; ok {
// 		rs.Mapped = val.(bool)
// 	}

// 	sql, err := args.Update("tboardgameprices")
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = db.NamedExec(sql.String(), &rs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return rs, nil
// }

// func DeletePrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
// 	rs, err := GetPriceById(db, args.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = db.NamedExec(`delete from tboardgameprices where id = :id`, &rs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return rs, nil
// }

// func UnmapPrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
// 	rs, err := GetPriceById(db, args.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = db.NamedExec(`update tboardgameprices set boardgame_id = NULL where id = :id`, &rs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return rs, nil
// }
