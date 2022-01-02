package boardgames

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/boardgames/search"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"log"
)

func GetSearch(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Price

	boardgame, err := getBoardgame(db, args.Id)
	if err != nil {
		return nil, err
	}

	if boardgame == nil {
		return nil, errors.New("Boardgame not found in the database")
	}

	tmp, err := search.Boardgame(*boardgame)
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

	for _, boardgame := range boardgames {
		log.Println(boardgame.Name)
		tmp, err := search.Boardgame(boardgame)
		if err != nil {
			return nil, err
		}
		rs = append(rs, tmp...)
	}

	return rs, nil
}
