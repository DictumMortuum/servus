package boardgames

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetDuel(c *gin.Context) {
	rs, err := getDuel()
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Success(c, &rs)
}

func getDuel() ([]models.Stats, error) {
	var rs []models.Stats

	database, err := db.Conn()
	if err != nil {
		return rs, err
	}
	defer database.Close()

	rs, err = getStats(database)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func getStats(db *sqlx.DB) ([]models.Stats, error) {
	var rs []models.Stats

	// sql := `
	// select
	//   s.*,
	// 	st.data,
	//   p.Name as player_name,
	//   IFNULL(exp.pantheon_points, 0) as pantheon_points,
	//   IF(exp.pantheon_points is not null, "pantheon", "") as expansions
	// from
	//   tduelstats s
	// join
	//   tboardgames g on g.name = "7 wonders duel"
	// join
	//   tboardgameplayers p on s.player_id = p.id
	// join
	//   tboardgameplays pl on s.play_id = pl.id
	// join
	// 	tboardgamestats st on s.play_id = st.play_id and s.player_id = st.player_id
	// order by play_id, s.id`

	sql := `
	select
		s.*
	from
		tboardgamestats s,
		tboardgameplays p,
		tboardgames g
	where
		g.name = "7 wonders duel" and
		g.id = p.boardgame_id and
		p.id = s.play_id
	`

	err := db.Select(&rs, sql)
	if err != nil {
		return rs, err
	}

	return rs, nil
}
