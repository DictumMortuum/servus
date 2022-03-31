package boardgames

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"text/template"
)

func GetPriceById(db *sqlx.DB, id int64) (*models.Price, error) {
	var rs models.Price

	err := db.QueryRowx(`
		select
			p.*,
			g.rank,
			IFNULL(g.thumb,"") as thumb,
			IFNULL(g.name,"") as boardgame_name,
			s.name as store_name
		from
			tboardgameprices p
			left join tboardgames g on g.id = p.boardgame_id,
			tboardgamestores s
		where
			p.id = ? and
			p.store_id = s.id
	`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func GetHistoricPricesById(db *sqlx.DB, id int64) ([]models.HistoricPrice, error) {
	var rs []models.HistoricPrice

	sql := `
		select
			boardgame_id,
			cr_date,
			avg(price) avg,
			min(price) min,
			max(price) max
		from
			tboardgamepriceshistory
		where
			boardgame_id = ?
		group by 1,YEAR(cr_date),MONTH(cr_date)
	`

	err := db.Select(&rs, sql, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetPrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return GetPriceById(db, args.Id)
}

func countPrices(db *sqlx.DB, args *models.QueryBuilder) (int, error) {
	var tpl bytes.Buffer

	sql := `
		select
			1
		from
			tboardgameprices p
			left join tboardgames g on g.id = p.boardgame_id,
			tboardgamestores s
		where
			p.store_id = s.id
		{{ if gt (len .FilterVal) 0 }}
			and p.{{ .FilterKey }} = {{ .FilterVal }}
		{{ end }}
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

func GetListPrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Price

	count, err := countPrices(db, args)
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", count))

	sql := `
		select
			p.*,
			g.rank,
			IFNULL(g.thumb,"") as thumb,
			IFNULL(g.name,"") as boardgame_name,
			s.name as store_name
		from
			tboardgameprices p
			left join tboardgames g on g.id = p.boardgame_id,
			tboardgamestores s
		where
			p.store_id = s.id
		{{ if gt (len .Ids) 0 }}
			and p.{{ .RefKey }} in (?)
		{{ else if gt (len .FilterVal) 0 }}
			and p.{{ .FilterKey }} = {{ .FilterVal }}
		{{ end }}
		{{ if gt (len .Sort) 0 }}
		order by {{ .Sort }} {{ .Order }}
		{{ else }}
		order by g.rank asc, p.cr_date asc
		{{ end }}
		{{ if eq (len .Range) 2 }}
		limit {{ index .Range 0 }}, {{ .Page }}
		{{ end }}`

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

	if args.Resources["history"] {
		retval := []models.Price{}

		for _, item := range rs {
			item.HistoricPrices, err = GetHistoricPricesById(db, item.BoardgameId.Int64)
			if err != nil {
				return nil, err
			}
			retval = append(retval, item)
		}

		return retval, nil
	} else {
		return rs, nil
	}
}

func CreatePrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs models.Price

	if val, ok := args.Data["boardgame_id"]; ok {
		rs.BoardgameId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}
	} else {
		return nil, errors.New("please provide a 'boardgame_id' parameter")
	}

	if val, ok := args.Data["store_id"]; ok {
		rs.StoreId = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'store_id' parameter")
	}

	if val, ok := args.Data["price"]; ok {
		rs.Price = val.(float64)
	} else {
		return nil, errors.New("please provide a 'price' parameter")
	}

	if val, ok := args.Data["stock"]; ok {
		rs.Stock = val.(bool)
	} else {
		return nil, errors.New("please provide a 'stock' parameter")
	}

	if val, ok := args.Data["url"]; ok {
		rs.Url = val.(string)
	} else {
		return nil, errors.New("please provide a 'url' parameter")
	}

	if val, ok := args.Data["store_thumb"]; ok {
		rs.StoreThumb = val.(string)
	} else {
		return nil, errors.New("please provide a 'store_thumb' parameter")
	}

	if val, ok := args.Data["batch"]; ok {
		rs.Batch = val.(int64)
	} else {
		return nil, errors.New("please provide a 'batch' parameter")
	}

	rs.Levenshtein = 0
	rs.Hamming = 0

	query, err := args.Insert("tboardgameprices")
	if err != nil {
		return nil, err
	}

	price, err := db.NamedExec(query.String(), &rs)
	if err != nil {
		return nil, err
	}

	id, err := price.LastInsertId()
	if err != nil {
		return nil, err
	}

	rs.Id = id
	return rs, nil
}

func UpdatePrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := GetPriceById(db, args.Id)
	if err != nil {
		return nil, err
	}

	args.IgnoreColumn("rank")
	args.IgnoreColumn("thumb")
	args.IgnoreColumn("boardgame_name")
	args.IgnoreColumn("store_name")

	if val, ok := args.Data["boardgame_id"]; ok {
		rs.BoardgameId = models.JsonNullInt64{
			Int64: int64(val.(float64)),
			Valid: true,
		}

		exists, err := boardgameExists(db, map[string]interface{}{
			"id": rs.BoardgameId.Int64,
		})
		if err != nil {
			return nil, err
		}

		if exists == nil {
			_, err := bgg.FetchBoardgame(db, rs.BoardgameId.Int64)
			if err != nil {
				return nil, err
			}
		}
	}

	if val, ok := args.Data["store_id"]; ok {
		rs.StoreId = int64(val.(float64))
	}

	if val, ok := args.Data["store_thumb"]; ok {
		rs.StoreThumb = val.(string)
	}

	if val, ok := args.Data["price"]; ok {
		rs.Price = val.(float64)
	}

	if val, ok := args.Data["stock"]; ok {
		rs.Stock = val.(bool)
	}

	if val, ok := args.Data["url"]; ok {
		rs.Url = val.(string)
	}

	if val, ok := args.Data["batch"]; ok {
		rs.Batch = int64(val.(float64))
	}

	if val, ok := args.Data["mapped"]; ok {
		rs.Mapped = val.(bool)
	}

	sql, err := args.Update("tboardgameprices")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func DeletePrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := GetPriceById(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgameprices where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func UnmapPrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := GetPriceById(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`update tboardgameprices set boardgame_id = NULL where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
