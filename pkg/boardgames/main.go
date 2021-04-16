package boardgames

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"log"
)

func GetDuel() ([]models.DuelPlay, error) {
	var rs []models.DuelPlay

	database, err := db.Conn()
	if err != nil {
		return rs, err
	}
	defer database.Close()

	rs, err = getDuels(database)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func GetWingspan() ([]models.WingspanPlay, error) {
	var rs []models.WingspanPlay

	database, err := db.Conn()
	if err != nil {
		return rs, err
	}
	defer database.Close()

	rs, err = getWingspan(database)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func GetWingspan2(id int64) (interface{}, error) {
	var rs []models.DuelPlay

	database, err := db.Conn()
	if err != nil {
		return rs, err
	}
	defer database.Close()

	log.Println(id)

	rs, err = getDuels(database)
	if err != nil {
		return rs, err
	}

	return rs, nil
}
