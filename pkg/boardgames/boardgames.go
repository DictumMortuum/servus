package boardgames

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"text/template"
	"time"
)

func getTopBoardgames(db *sqlx.DB) ([]models.Boardgame, error) {
	var rs []models.Boardgame

	sql := `
		select
			*
		from
			tboardgames
		where
			rank <= 100
		order by rank
	`

	err := db.Select(&rs, sql)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

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

func boardgameExists(db *sqlx.DB, payload map[string]interface{}) (*models.JsonNullInt64, error) {
	var id models.JsonNullInt64

	q := `select id from tboardgames where id = :id`
	stmt, err := db.PrepareNamed(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Get(&id, payload)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func GetBoardgame(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return getBoardgame(db, args.Id)
}

func countBoardgames(db *sqlx.DB, args *models.QueryBuilder) (int, error) {
	var tpl bytes.Buffer

	sql := `
		select
			1
		from
			tboardgames b
		{{ if gt (len .Ids) 0 }}
		where b.{{ .RefKey }} in (?)
		{{ else if eq .FilterKey "ranked"}}
		where b.rank is not null
		{{ else if gt (len .FilterVal) 0 }}
		where b.{{ .FilterKey }} = "{{ .FilterVal }}"
		{{ end }}
		{{ if gt (len .Query) 0 }}
		and b.name like "%{{ .Query }}%"
		{{ end }}
	`

	t := template.Must(template.New("count").Parse(sql))
	err := t.Execute(&tpl, args)
	if err != nil {
		return -1, err
	}

	query, ids, err := sqlx.In(tpl.String(), args.Ids)
	if err != nil {
		query = tpl.String()
	} else {
		query = db.Rebind(query)
	}

	var count []int
	err = db.Select(&count, query, ids...)
	if err != nil {
		return -1, err
	}

	return len(count), nil
}

func GetListBoardgame(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Boardgame

	count, err := countBoardgames(db, args)
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", count))

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
		{{ else if eq .FilterKey "ranked" }}
		where b.rank is not null
		{{ else if gt (len .FilterVal) 0 }}
		where b.{{ .FilterKey }} = "{{ .FilterVal }}"
		{{ end }}
		{{ if gt (len .Query) 0 }}
		and b.name like "%{{ .Query }}%"
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

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		err := fn(args.Data, true)
		if err != nil {
			return nil, err
		}
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

func RefetchBoardgame(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	_, err := bgg.FetchBoardgame(db, args.Id)
	if err != nil {
		return nil, err
	}

	rs, err := getBoardgame(db, args.Id)
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

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		fn(args.Data, false)
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
