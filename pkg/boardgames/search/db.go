package search

import (
	"database/sql"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
)

func upsertPrice(db *sqlx.DB, item models.Price) error {
	item.BoardgameId = models.JsonNullInt64{
		Int64: -1,
		Valid: false,
	}

	id, err := findPrice(db, item)
	if err != nil {
		return err
	}

	if id == nil {
		_, err := create(db, item)
		if err != nil {
			return err
		}
	} else {
		item.Id = id.Int64
		_, err := update(db, item)
		if err != nil {
			return err
		}
	}

	return nil
}

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

func findPrice(db *sqlx.DB, payload models.Price) (*models.JsonNullInt64, error) {
	var id models.JsonNullInt64

	q := `
		select
			id
		from
			tboardgameprices
		where
			store_id = :store_id and
			name = :name
	`

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

func exists(db *sqlx.DB, payload models.Price) (bool, error) {
	q := `select 1 from tboardgameprices where store_id = :store_id and name = :name`

	rows, err := db.NamedQuery(q, payload)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func updateBatch(db *sqlx.DB, store_id int64) error {
	q := `
		update
			tboardgameprices
		set
			batch = 0
		where
			store_id = :store_id
	`

	_, err := db.NamedExec(q, map[string]interface{}{
		"store_id": store_id,
	})
	if err != nil {
		return err
	}

	return nil
}

func update(db *sqlx.DB, payload models.Price) (bool, error) {
	q := `
		update
			tboardgameprices
		set
			store_thumb = :store_thumb,
			price = :price,
			stock = :stock,
			url = :url,
			batch = 1,
			cr_date = NOW()
		where
			id = :id
	`

	rs, err := db.NamedExec(q, payload)
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func create(db *sqlx.DB, payload models.Price) (bool, error) {
	q := `
		insert into tboardgameprices (
			boardgame_id,
			name,
			store_id,
			store_thumb,
			price,
			stock,
			url,
			levenshtein,
			hamming,
			batch
		) values (
			:boardgame_id,
			:name,
			:store_id,
			:store_thumb,
			:price,
			:stock,
			:url,
			:levenshtein,
			:hamming,
			1
		)`

	rs, err := db.NamedExec(q, payload)
	if err != nil {
		return false, err
	}

	rows, err := rs.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
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