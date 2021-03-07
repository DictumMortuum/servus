package boardgames

import (
	"github.com/DictumMortuum/servus/db"
)

func GetDuel() ([]DuelStatsRow, error) {
	var rs []DuelStatsRow

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

func GetWingspan() ([]WingspanStatsRow, error) {
	var rs []WingspanStatsRow

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
