package boardgames

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"log"
	"reflect"
	"time"
)

func getWingspan(db *sqlx.DB) ([]models.WingspanPlay, error) {
	var rs []models.WingspanPlay

	sql := `
  select
    s.*,
    p.Name as player_name
  from
    twingspanstats s
  join
    tboardgames g on g.name = "wingspan"
  join
    tboardgameplayers p on s.player_id = p.id
  join
    tboardgameplays pl on s.play_id = pl.id
	order by play_id, s.id`

	err := db.Select(&rs, sql)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func stringToDateTimeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
		return time.Parse(time.RFC3339, data.(string))
	}

	return data, nil
}

func insertWingspan(db *sqlx.DB, data map[string]interface{}) (int64, error) {
	var obj models.WingspanPlay

	log.Println("1", data)

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     &obj,
		TagName:    "json",
		DecodeHook: stringToDateTimeHook,
	})

	err = decoder.Decode(data)
	if err != nil {
		return -1, err
	}

	log.Println("asdf", obj)

	return -1, nil
	/*
		player_id, err := upsertPlayer(db, obj.PlayerName)
		if err != nil {
			return -1, err
		}

		obj.PlayerId = player_id

		boardgame, err := getBoardgame(db, "wingspan")
		if err != nil {
			return -1, err
		}

		date, err := time.Parse("2006-01-02 15:04:05", data["date"].(string))
		if err != nil {
			return -1, err
		}

		play_id, err := createPlay(db, PlayModel{
			Date:        date,
			BoardgameId: boardgame.Id,
		})
	*/
	obj.PlayId = 1

	rs, err := db.NamedExec(`
	insert into twingspanstats (
		play_id,
		player_id,
		bird_points,
		bonus_points,
		endofround_points,
		egg_points,
		food_points,
		tucked_points
	) values (
		:play_id,
		:player_id,
		:bird_points,
		:bonus_points,
		:endofround_points,
		:egg_points,
		:food_points,
		:tucked_points
	)`, &obj)
	if err != nil {
		return -1, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
