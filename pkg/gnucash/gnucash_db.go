package gnucash

import (
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

type ExpensesRow struct {
	Date time.Time `db:"start_of_month"`
	Name string    `db:"name"`
	Sum  float64   `db:"sum"`
}

func getExpenseByMonthWithDateFilter(db *sqlx.DB, name string, raw_from string, raw_to string) ([]ExpensesRow, error) {
	rs := []ExpensesRow{}

	from, err := strconv.ParseInt(raw_from, 10, 64)
	if err != nil {
		return nil, err
	}

	to, err := strconv.ParseInt(raw_to, 10, 64)
	if err != nil {
		return nil, err
	}

	from_date := time.UnixMilli(from)
	to_date := time.UnixMilli(to)
	names := strings.Split(name, ",")

	sql := `
		select a.start_of_month, name, sum(a.sum) as sum from
		(
			select
				DATE_SUB(t.post_date,INTERVAL DAYOFMONTH(t.post_date)-1 DAY) as start_of_month,
				"placeholder" as name,
				0 as sum
			from
				transactions t
			union
			select
				DATE_SUB(t.post_date,INTERVAL DAYOFMONTH(t.post_date)-1 DAY) as start_of_month,
				a.name,
				sum(1.0*s.value_num/s.value_denom) as sum
			from
				transactions t,
				splits s,
				accounts a
			where
				t.guid = s.tx_guid and
				a.guid = s.account_guid and
				a.name in (:name)
			group by 1, 2
			order by 1
		) a
		where
			a.start_of_month > :from_date and a.start_of_month < :to_date
		group by 1, 2
		order by 1
	`

	arg := map[string]interface{}{
		"from_date": from_date,
		"to_date":   to_date,
		"name":      names,
	}

	query, ids, err := sqlx.Named(sql, arg)
	query, ids, err = sqlx.In(query, ids...)
	query = db.Rebind(query)
	err = db.Select(&rs, query, ids...)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func getExpenseByMonth(db *sqlx.DB, name string) ([]ExpensesRow, error) {
	rs := []ExpensesRow{}

	names := strings.Split(name, ",")

	sql := `
		select a.start_of_month, name, sum(a.sum) as sum from
		(
			select
				DATE_SUB(t.post_date,INTERVAL DAYOFMONTH(t.post_date)-1 DAY) as start_of_month,
				"placeholder" as name,
				0 as sum
			from
				transactions t
			union
			select
				DATE_SUB(t.post_date,INTERVAL DAYOFMONTH(t.post_date)-1 DAY) as start_of_month,
				a.name,
				sum(1.0*s.value_num/s.value_denom) as sum
			from
				transactions t,
				splits s,
				accounts a
			where
				t.guid = s.tx_guid and
				a.guid = s.account_guid and
				a.name in (?)
			group by 1, 2
			order by 1
		) a
		group by 1, 2
		order by 1
	`

	query, ids, err := sqlx.In(sql, names)
	if err != nil {
		query = sql
	} else {
		query = db.Rebind(query)
	}

	err = db.Select(&rs, query, ids...)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func getTotalExpenseByMonth(db *sqlx.DB) ([]ExpensesRow, error) {
	rs := []ExpensesRow{}

	err := db.Select(&rs, `
	select
		a.name
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
		a.code != "tax"
	group by 1,2
	order by 1,2
	`)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

type TopExpensesRow struct {
	Name    string  `db:"name" json:"name"`
	Total   float64 `db:"total" json:"total"`
	Average float64 `db:"average" json:"average"`
}

func getTopExpenses(db *sqlx.DB, date string) ([]TopExpensesRow, error) {
	rs := []TopExpensesRow{}

	err := db.Select(&rs, `
	select
		a.name,
		TRUNCATE(sum(1.0*s.value_num/s.value_denom), 2) as total,
		TRUNCATE(sum(1.0*s.value_num/s.value_denom)/TIMESTAMPDIFF(MONTH, ?, NOW()), 2) as average
	from
		transactions t,
		splits s,
		accounts a
	where
		t.guid = s.tx_guid and
		a.guid = s.account_guid and
		a.account_type= "EXPENSE" and
		t.post_date > ? and
		a.code != "tax"
	group by 1
	order by 2 desc, 1
	limit 15
	`, date, date)
	if err != nil {
		return rs, err
	}

	return rs, nil
}
