package models

type Stats struct {
	Id            int64   `db:"id" json:"id"`
	PlayId        int64   `db:"play_id" json:"play_id"`
	BoardgameId   int64   `db:"boardgame_id" json:"boardgame_id"`
	PlayerId      int64   `db:"player_id" json:"player_id"`
	Data          Json    `db:"data" json:"data"`
	Player        string  `db:"name" json:"player"`
	PlayerSurname string  `db:"surname" json:"player_surname"`
	Mu            float64 `json:"mu"`
	Sigma         float64 `json:"sigma"`
	TrueSkill     float64 `json:"trueskill"`
	Delta         float64 `json:"delta"`
}
