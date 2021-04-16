package boardgames

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

func getDuels(db *sqlx.DB) ([]models.DuelPlay, error) {
	var rs []models.DuelPlay

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

func createDuelStats(db *sqlx.DB, data models.DuelPlay) error {
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
