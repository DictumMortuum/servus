package gnucash

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type ExpensesRow struct {
	Date time.Time `db:"start_of_month"`
	Sum  float64   `db:"sum"`
}

func getExpenseByMonth(db *sqlx.DB, name string) ([]ExpensesRow, error) {
	rs := []ExpensesRow{}

	err := db.Select(&rs, `
	select
		DATE_SUB(t.post_date,INTERVAL DAYOFMONTH(t.post_date)-1 DAY) as start_of_month,
		sum(1.0*s.value_num/s.value_denom) as sum
	from
		transactions t,
		splits s,
		accounts a
	where
		t.guid = s.tx_guid and
		a.guid = s.account_guid and
		a.account_type= "EXPENSE" and
		a.name = ?
	group by 1
	order by 1
	`, name)
	if err != nil {
		return rs, err
	}

	return rs, nil
}
