package boardgames

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"text/template"
	"time"
)

func getBoardgame(db *sqlx.DB, id int64) (*models.Boardgame, error) {
	var rs models.Boardgame

	sql := `
		select
			b.*,
			TRUNCATE(sum(1.0*s.value_num/s.value_denom), 2) as cost
		from
			tboardgames b
			left join gnucash.transactions t on t.guid = b.tx_guid
			left join gnucash.splits s on s.tx_guid = b.tx_guid and s.account_guid = "3097dd8d65751277845bdda438cba937"
		where
			b.id = ?
		group by 1
	`

	err := db.QueryRowx(sql, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	rs.Date = time.Now().AddDate(1, 0, 0)

	return &rs, nil
}

func GetBoardgame(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getBoardgame(db, args.Id)
}

func GetListBoardgame(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Boardgame

	var count []int
	err := db.Select(&count, "select 1 from tboardgames")
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", len(count)))

	sql := `
		select
			b.*,
			TRUNCATE(sum(1.0*s.value_num/s.value_denom), 2) as cost
		from
			tboardgames b
			left join gnucash.transactions t on t.guid = b.tx_guid
			left join gnucash.splits s on s.tx_guid = b.tx_guid and s.account_guid = "3097dd8d65751277845bdda438cba937"
		{{ if gt (len .Ids) 0 }}
		where b.{{ .RefKey }} in (?)
		{{ else if gt (len .FilterVal) 0 }}
		where b.{{ .FilterKey }} = "{{ .FilterVal }}"
		{{ end }}
		group by 1
		{{ if gt (len .Sort) 0 }}
		order by b.{{ .Sort }} {{ .Order }}
		{{ else }}
		order by b.id
		{{ end }}
		{{ if eq (len .Range) 2 }}
		limit {{ index .Range 0 }}, {{ .Page }}
		{{ end }}
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

func CreateBoardgame(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs models.Boardgame

	if val, ok := args.Data["id"]; ok {
		rs.Id = int64(val.(float64))
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	if val, ok := args.Data["name"]; ok {
		rs.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	if val, ok := args.Data["data"]; ok {
		err := rs.Data.Scan(val)
		if err != nil {
			return nil, err
		}
	} else {
		rs.Data = nil
	}

	query, err := args.Insert("tboardgames")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(query.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func UpdateBoardgame(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getBoardgame(db, args.Id)
	if err != nil {
		return nil, err
	}

	if val, ok := args.Data["name"]; ok {
		rs.Name = val.(string)
	}

	if val, ok := args.Data["data"]; ok {
		err := rs.Data.Scan(val)
		if err != nil {
			return nil, err
		}
	}

	sql, err := args.Update("tboardgames")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func DeleteBoardgame(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := getBoardgame(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgames where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
