package boardgames

import (
	"database/sql"
	"errors"
	"github.com/DictumMortuum/servus/pkg/boardgames/search"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
)

func getNextBatch(db *sqlx.DB) (*models.JsonNullInt64, error) {
	var id models.JsonNullInt64

	err := db.Get(&id, `select max(batch)+1 as next_batch from tboardgameprices`)
	if err == sql.ErrNoRows {
		return nil, errors.New("Could not find next batch_id")
	}
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func GetSearch(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Price

	boardgame, err := getBoardgame(db, args.Id)
	if err != nil {
		return nil, err
	}

	if boardgame == nil {
		return nil, errors.New("Boardgame not found in the database")
	}

	batch_id, err := getNextBatch(db)
	if err != nil {
		return nil, err
	}

	tmp, err := search.Boardgame(*boardgame, *batch_id)
	if err != nil {
		return nil, err
	}
	rs = append(rs, tmp...)

	return rs, nil
}

func SearchTop(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Price

	boardgames, err := getTopBoardgames(db)
	if err != nil {
		return nil, err
	}

	batch_id, err := getNextBatch(db)
	if err != nil {
		return nil, err
	}

	for _, boardgame := range boardgames {
		log.Println(boardgame.Rank, boardgame.Name)
		tmp, err := search.Boardgame(boardgame, *batch_id)
		if err != nil {
			return nil, err
		}
		rs = append(rs, tmp...)
	}

	return rs, nil
}
