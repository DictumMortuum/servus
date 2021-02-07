package boardgames

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type DuelRow struct {
	Play    DuelPlaysRow   `json:"play"`
	Players []DuelStatsRow `json:"players"`
}

type DuelPlaysRow struct {
	Id         int64     `db:"id" json:"id"`
	CrDate     time.Time `db:"cr_date" json:"cr_date"`
	Date       time.Time `db:"date" json:"date"`
	Expansions string    `db:"expansions" json:"expansions"`
}

type DuelStatsRow struct {
	Id             int64  `db:"id" json:"id"`
	PlayId         int64  `db:"play_id"`
	PlayerId       int64  `db:"player_id"`
	PlayerName     string `json:"player"`
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

func createDuelPlay(db *sqlx.DB, data DuelPlaysRow) (int64, error) {
	res, err := db.NamedExec(`
	insert into tduelplays (
		cr_date,
		date
	) values (
		NOW(),
		:date
	)`, &data)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func CreateDuelStats(db *sqlx.DB, data DuelStatsRow) error {
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

func getDuel(db *sqlx.DB, id int64) (*DuelRow, error) {
	var retval DuelRow
	var stats []DuelStatsRow

	err := db.QueryRowx(`
	select
		dp.*,
		IF(dpe.pantheon_points is not null, "pantheon", "") as expansions
	from
		tduelplays dp
	join
		tduelstats ds on dp.id = ds.play_id
	left join
		tduelpantheonexpansion dpe on dpe.stats_id = ds.id
	where
		dp.id = ?`, id).StructScan(&retval.Play)
	if err != nil {
		return nil, err
	}

	err = db.Select(&stats, `
	select
		ds.*,
		IFNULL(dpe.pantheon_points, 0) as pantheon_points
	from
		tduelstats ds
	left join
		tduelpantheonexpansion dpe on dpe.stats_id = ds.id
	where
		ds.play_id = ?`, id)
	if err != nil {
		return nil, err
	}

	for _, stat := range stats {
		player, err := GetPlayer(db, stat.PlayerId)
		if err != nil {
			return nil, err
		}

		stat.PlayerName = player.Name
		retval.Players = append(retval.Players, stat)
	}

	return &retval, nil
}

func getDuels(db *sqlx.DB) ([]DuelRow, error) {
	retval := []DuelRow{}
	var plays []DuelPlaysRow

	err := db.Select(&plays, `select * from tduelplays`)
	if err != nil {
		return retval, err
	}

	for _, play := range plays {
		duel, err := getDuel(db, play.Id)
		if err != nil {
			return retval, err
		}

		retval = append(retval, *duel)
	}

	return retval, nil
}
