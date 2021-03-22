package boardgames

import (
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"log"
	"reflect"
	"time"
)

type WingspanModel struct {
	Id               int64     `db:"id" json:"id"`
	PlayId           int64     `db:"play_id" json:"play_id"`
	PlayerId         int64     `db:"player_id" json:"player_id"`
	Date             time.Time `json:"date"`
	PlayerName       string    `db:"player_name" json:"player"`
	BirdPoints       int       `db:"bird_points" json:"birds"`
	BonusPoints      int       `db:"bonus_points" json:"bonus"`
	EndofroundPoints int       `db:"endofround_points" json:"endofround"`
	EggPoints        int       `db:"egg_points" json:"egg"`
	FoodPoints       int       `db:"food_points" json:"food"`
	TuckedPoints     int       `db:"tucked_points" json:"tucked"`
}

func getWingspan(db *sqlx.DB) ([]WingspanModel, error) {
	var rs []WingspanModel

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
	var obj WingspanModel

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
