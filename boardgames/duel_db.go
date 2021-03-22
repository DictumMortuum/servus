package boardgames

import (
	"github.com/jmoiron/sqlx"
)

type DuelModel struct {
	Id             int64  `db:"id" json:"id"`
	PlayId         int64  `db:"play_id" json:"play_id"`
	PlayerId       int64  `db:"player_id" json:"player_id"`
	PlayerName     string `db:"player_name" json:"player"`
	Expansions     string `db:"expansions" json:"expansions"`
	BluePoints     int    `db:"blue_points" json:"blue"`
	GreenPoints    int    `db:"green_points" json:"green"`
	YellowPoints   int    `db:"yellow_points" json:"yellow"`
	PurplePoints   int    `db:"purple_points" json:"purple"`
	WonderPoints   int    `db:"wonder_points" json:"wonder"`
	MarkerPoints   int    `db:"marker_points" json:"marker"`
	CoinPoints     int    `db:"coin_points" json:"coin"`
	BattlePoints   int    `db:"battle_points" json:"battle"`
	PantheonPoints int    `db:"pantheon_points" json:"pantheon"`
	BattleVictory  bool   `db:"battle_victory" json:"battle_victory"`
	ScienceVictory bool   `db:"science_victory" json:"science_victory"`
}

func getDuels(db *sqlx.DB) ([]DuelModel, error) {
	var rs []DuelModel

	sql := `
  select
    s.*,
    p.Name as player_name,
    IFNULL(exp.pantheon_points, 0) as pantheon_points,
    IF(exp.pantheon_points is not null, "pantheon", "") as expansions
  from
    tduelstats s
  join
    tboardgames g on g.name = "7 wonders duel"
  join
    tboardgameplayers p on s.player_id = p.id
  join
    tboardgameplays pl on s.play_id = pl.id
  left join
    tduelpantheonexpansion exp on exp.stats_id = s.id
	order by play_id, s.id`

	err := db.Select(&rs, sql)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func createDuelStats(db *sqlx.DB, data DuelModel) error {
	_, err := db.NamedExec(`
	insert into tduelstats (
		play_id,
		player_id,
		blue_points,
		green_points,
		yellow_points,
		purple_points,
		wonder_points,
		marker_points,
		coin_points,
		battle_points,
		battle_victory,
		science_victory
	) values (
		:play_id,
		:player_id,
		:blue_points,
		:green_points,
		:yellow_points,
		:purple_points,
		:wonder_points,
		:marker_points,
		:coin_points,
		:battle_points,
		:battle_victory,
		:science_victory
	)`, &data)
	if err != nil {
		return err
	}

	return nil
}
