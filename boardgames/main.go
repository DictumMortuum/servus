package boardgames

import (
	"github.com/DictumMortuum/servus/db"
	"log"
)

func GetDuel() ([]DuelModel, error) {
	var rs []DuelModel

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

func GetWingspan() ([]WingspanModel, error) {
	var rs []WingspanModel

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
	var rs []DuelModel

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
