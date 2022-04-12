package search

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
)

func normalize_dates(db *sqlx.DB) (bool, error) {
	q := `
		update tboardgamepriceshistory set cr_date = date_add(date_add(LAST_DAY(cr_date), interval 1 day), interval -1 month)
	`

	rs, err := db.NamedExec(q, map[string]interface{}{})
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func insert_mappings(db *sqlx.DB) (bool, error) {
	q := `
		insert into tboardgamepricesmap (
			boardgame_id,
			name
		)
		select distinct
			p.boardgame_id,
			p.name
		from
			tboardgameprices p
		where
			boardgame_id is not null and
			not exists (select 1 from tboardgamepricesmap where name = p.name);
	`

	rs, err := db.NamedExec(q, map[string]interface{}{})
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func insert_histories(db *sqlx.DB) (bool, error) {
	q := `
		insert into	tboardgamepriceshistory (
			boardgame_id,
			cr_date,
			price,
			stock,
			store_id
		)
		select
			boardgame_id,
			cr_date,
			price,
			stock,
			store_id
		from
			tboardgameprices
		where
			boardgame_id is not null and
			mapped = 1 and
			batch = 1
	`

	rs, err := db.NamedExec(q, map[string]interface{}{})
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func update_mapping(db *sqlx.DB) (bool, error) {
	q := `
		update tboardgameprices set mapped = 1 where boardgame_id is not null
	`

	rs, err := db.NamedExec(q, map[string]interface{}{})
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func delete_redundant_prices(db *sqlx.DB) (int64, error) {
	q := `
		delete from
			tboardgameprices
		where
			id in (select
							id
						from
							(select
								name,
								store_id,
								max(id) as id
							from
								tboardgameprices
							group by 1,2
							having count(*) > 1)
						p);
	`

	rs, err := db.NamedExec(q, map[string]interface{}{})
	if err != nil {
		return -1, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

func delete_redundant_history(db *sqlx.DB) (int64, error) {
	q := `
		delete from
			tboardgamepriceshistory
		where
			id in (select
							id
						from
							(select
								boardgame_id,
								price,
								cr_date,
								store_id,
								min(id) as id
							from
								tboardgamepriceshistory
							group by 1,2,3,4
							having count(*) > 1)
						p);
	`

	rs, err := db.NamedExec(q, map[string]interface{}{})
	if err != nil {
		return -1, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

func UpdateMappings(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	_, err := update_mapping(db)
	if err != nil {
		return nil, err
	}

	_, err = insert_mappings(db)
	if err != nil {
		return nil, err
	}

	_, err = insert_histories(db)
	if err != nil {
		return nil, err
	}

	_, err = normalize_dates(db)
	if err != nil {
		return nil, err
	}

	var count int64
	count = 1
	for count > 0 {
		count, err = delete_redundant_prices(db)
		log.Println("Found:", count)
		if err != nil {
			return nil, err
		}
	}

	count = 1
	for count > 0 {
		count, err = delete_redundant_history(db)
		log.Println("Found:", count)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
