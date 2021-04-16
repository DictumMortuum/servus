package gnucash

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"time"
)

type ExpensesRow struct {
	RawMonth string  `db:"month"`
	Sum      float64 `db:"sum"`
}

func (e ExpensesRow) Month() (*time.Time, error) {
	time, err := time.Parse("2006-01", e.RawMonth)
	if err != nil {
		return nil, err
	}
	return &time, nil
}

func (e ExpensesRow) MarshalJSON() ([]byte, error) {
	month, _ := e.Month()

	return json.Marshal(struct {
		Month time.Time `json:"month"`
		Sum   float64   `json:"sum"`
	}{
		Month: *month,
		Sum:   e.Sum,
	})
}

func getExpenseByMonth(db *sqlx.DB, name string) ([]ExpensesRow, error) {
	rs := []ExpensesRow{}

	err := db.Select(&rs, `
	select
		DATE_FORMAT(t.post_date, '%Y-%m') as month,
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
