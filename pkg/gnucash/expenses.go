package gnucash

import (
	"bytes"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"text/template"
	"time"
)

type Expense struct {
	Date     time.Time `db:"post_date" json:"date"`
	Desc     string    `db:"description" json:"description"`
	Guid     string    `db:"guid" json:"id"`
	Tx_Guid  string    `db:"tx_guid" json:"tx_guid"`
	Acc_Guid string    `db:"account_guid" json:"account_guid"`
	Price    float64   `db:"price" json:"price"`
	Name     string    `db:"name" json:"name"`
	Type     string    `db:"account_type" json:"type"`
}

func GetListExpenses(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []Expense

	var count []int
	err := db.Select(&count, "select 1 from splits")
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", len(count)))

	sql := `
		select
			t.post_date,
			t.description,
			s.guid,
			s.tx_guid,
			s.account_guid,
			1.0*s.value_num/s.value_denom as price,
			a.name,
			a.account_type
		from
			transactions t,
			splits s,
			accounts a
		where
			t.guid = s.tx_guid and
			a.guid = s.account_guid
		order by 1
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
