package boardgames

import (
	"github.com/jmoiron/sqlx"
)

type WingspanStatsRow struct {
	Id               int64  `db:"id" json:"id"`
	PlayId           int64  `db:"play_id" json:"play_id"`
	PlayerId         int64  `db:"player_id" json:"player_id"`
	PlayerName       string `db:"player_name" json:"player"`
	BirdPoints       int    `db:"bird_points" json:"birds"`
	BonusPoints      int    `db:"bonus_points" json:"bonus"`
	EndofroundPoints int    `db:"endofround_points" json:"endofround"`
	EggPoints        int    `db:"egg_points" json:"egg"`
	FoodPoints       int    `db:"food_points" json:"food"`
	TuckedPoints     int    `db:"tucked_points" json:"tucked"`
}

func getWingspan(db *sqlx.DB) ([]WingspanStatsRow, error) {
	var rs []WingspanStatsRow

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
