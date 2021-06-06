package boardgames

import (
	DB "github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"sort"
)

func GetScores(c *gin.Context) {
	type retval struct {
		Play  models.Play    `json:"play"`
		Stats []models.Stats `json:"stats"`
	}

	db, err := DB.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer db.Close()

	var plays []models.Play
	err = db.Select(&plays, "select * from tboardgameplays")
	if err != nil {
		util.Error(c, err)
		return
	}

	rs := []retval{}
	for _, play := range plays {
		f, g := getFuncs(play.BoardgameId)
		if f == nil || g == nil {
			continue
		}

		var stats []models.Stats
		err = db.Select(&stats, "select * from tboardgamestats where play_id = ?", play.Id)
		if err != nil {
			util.Error(c, err)
			return
		}

		retval2 := []models.Stats{}
		for _, item := range stats {
			item.Data["score"] = f(item)
			retval2 = append(retval2, item)
		}

		sort.Slice(retval2, g(retval2))
		rs = append(rs, retval{play, retval2})
	}

	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Play.Date.Before(rs[j].Play.Date)
	})

	util.Success(c, rs)
}

// MariaDB [servus]> select id, name from tboardgames where id in (select distinct boardgame_id from tboardgameplays);
// +--------+--------------------------------+
// | id     | name                           |
// +--------+--------------------------------+
// |     42 | Tigris & Euphrates             |
// |    822 | Carcassonne                    |
// |  14996 | Ticket to ride: Europe         |
// | 110327 | Lords of Waterdeep             |
// | 127023 | Kemet                          |
// | 163412 | Patchwork                      |
// | 170042 | Raiders of the north sea       |
// | 170216 | Blood rage                     |
// | 173346 | 7 Wonders Duel                 |
// | 183394 | Viticulture essential edition  |
// | 230802 | Azul                           |
// | 236457 | Architects of the west kingdom |
// | 237182 | Root                           |
// | 256916 | Concordia Venus                |
// | 266192 | Wingspan                       |
// | 266810 | Paladins of the West Kingdom   |
// | 271320 | Castles of Burgundy            |
// | 283863 | The Magnificent                |
// | 312484 | Lost ruins of Arnak            |
// +--------+--------------------------------+
// 19 rows in set (0.002 sec)

func getFuncs(boardgame_id int64) (func(models.Stats) float64, func([]models.Stats) func(i, j int) bool) {
	switch boardgame_id {
	case 183394:
		return DefaultScore, DefaultSort
	case 110327:
		return WaterdeepScore, WaterdeepSort
	case 173346:
		return DuelScore, DuelSort
	case 266192:
		return WingspanScore, WingspanSort
	case 266810:
		return PaladinsScore, PaladinsSort
	case 230802:
		return DefaultScore, AzulSort
	case 237182:
		return PlaceScore, RootSort
	case 283863:
		return DefaultScore, MagnificentSort
	case 312484:
		return LostRuinsScore, LostRuinsSort
	case 256916:
		return ConcordiaScore, ConcordiaSort
	case 163412:
		return PatchworkScore, PatchworkSort
	case 822:
		return DefaultScore, DefaultSort
	case 42:
		return DefaultScore, DefaultSort
	case 170042:
		return RaidersScore, RaidersSort
	case 14996:
		return DefaultScore, DefaultSort
	case 236457:
		return ArchitectsScore, ArchitectsSort
	case 170216:
		return DefaultScore, DefaultSort
	case 127023:
		return KemetScore, KemetSort
	case 271320:
		return DefaultScore, DefaultSort
	default:
		return nil, nil
	}
}

func DefaultScore(stats models.Stats) float64 {
	return stats.Data["score"].(float64)
}

func PlaceScore(stats models.Stats) float64 {
	return stats.Data["place"].(float64)
}

func DefaultSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		return DefaultScore(stats[i]) < DefaultScore(stats[j])
	}
}

// {
// 	"quests": 53,
// 	"diamonds": 5,
// 	"cubes": 8,
// 	"money": 4,
// 	"character": 8
// }

func WaterdeepScore(stats models.Stats) float64 {
	keys := []string{"quests", "diamonds", "cubes", "money", "character"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	return score
}

func WaterdeepSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := WaterdeepScore(stats[i])
		score2 := WaterdeepScore(stats[j])

		if int64(score1) == int64(score2) {
			return stats[i].Data["money"].(float64) < stats[j].Data["money"].(float64)
		} else {
			return score1 < score2
		}
	}
}

// {
// 	"battle_points":0,
// 	"battle_victory":false,
// 	"blue_points":10,
// 	"coin_points":5,
// 	"expansions":"",
// 	"green_points":6,
// 	"marker_points":13,
// 	"pantheon_points":0,
// 	"purple_points":0,
// 	"science_victory":false,
// 	"wonder_points":16,
// 	"yellow_points":9
// }

func DuelScore(stats models.Stats) float64 {
	keys := []string{
		"battle_points",
		"blue_points",
		"coin_points",
		"green_points",
		"marker_points",
		"pantheon_points",
		"purple_points",
		"wonder_points",
		"yellow_points",
	}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	return score
}

func DuelSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		if stats[i].Data["battle_victory"].(bool) {
			return false
		}

		if stats[j].Data["battle_victory"].(bool) {
			return true
		}

		if stats[i].Data["science_victory"].(bool) {
			return false
		}

		if stats[j].Data["science_victory"].(bool) {
			return true
		}

		score1 := DuelScore(stats[i])
		score2 := DuelScore(stats[j])

		if int64(score1) == int64(score2) {
			return stats[i].Data["coin_points"].(float64) < stats[j].Data["coin_points"].(float64)
		} else {
			return score1 < score2
		}
	}
}

// {
// 	"birds": 37,
// 	"bonus": 6,
// 	"endofround": 16,
// 	"eggs": 9,
// 	"food": 11,
// 	"tucked": 15
// }

func WingspanScore(stats models.Stats) float64 {
	keys := []string{"birds", "bonus", "endofround", "eggs", "food", "tucked"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	return score
}

func WingspanSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := WingspanScore(stats[i])
		score2 := WingspanScore(stats[j])
		return score1 < score2
	}
}

// {
// 	"red": 3,
// 	"black": 9,
// 	"blue": 11,
// 	"commission": 9,
// 	"fortify": 3,
// 	"garrison": 0,
// 	"absolve": 6,
// 	"attack": 0,
// 	"convert": 0,
// 	"order": 8,
// 	"develop": 6,
// 	"debt": 2,
// 	"coin": 0
// }

func PaladinsScore(stats models.Stats) float64 {
	keys := []string{"red", "black", "blue", "commission", "fortify", "garrison", "absolve", "attack", "convert", "order", "develop", "debt", "coin"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	return score
}

func PaladinsSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := PaladinsScore(stats[i])
		score2 := PaladinsScore(stats[j])

		if int64(score1) == int64(score2) {
			return stats[i].Data["order"].(float64) < stats[j].Data["order"].(float64)
		}

		return score1 < score2
	}
}

func AzulSort(stats []models.Stats) func(i, j int) bool {
	// {
	// 	"score": 108,
	// 	"colors": 1,
	// 	"rows": 1,
	// 	"columns": 0
	// }
	return func(i, j int) bool {
		score1 := DefaultScore(stats[i])
		score2 := DefaultScore(stats[j])

		if int64(score1) == int64(score2) {
			return stats[i].Data["rows"].(float64) < stats[j].Data["rows"].(float64)
		}

		return score1 < score2
	}
}

func MagnificentSort(stats []models.Stats) func(i, j int) bool {
	// {
	// 	"player": 1,
	// 	"score": 99
	// }
	return func(i, j int) bool {
		return DefaultScore(stats[i]) < DefaultScore(stats[j])
	}
}

func RootSort(stats []models.Stats) func(i, j int) bool {
	// {
	// 	"faction":"marquise de cat",
	// 	"place":1
	// }
	return func(i, j int) bool {
		return PlaceScore(stats[i]) > PlaceScore(stats[j])
	}
}

// {
//   "research": 29,
//   "tiles": 4,
//   "idols": 12,
//   "guardian": 20,
//   "cards": 6,
//   "fear": -1
// }

func LostRuinsScore(stats models.Stats) float64 {
	keys := []string{"research", "tiles", "idols", "guardian", "cards", "fear"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	return score
}

func LostRuinsSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := LostRuinsScore(stats[i])
		score2 := LostRuinsScore(stats[j])

		if int64(score1) == int64(score2) {
			return stats[i].Data["research"].(float64) < stats[j].Data["research"].(float64)
		}

		return score1 < score2
	}
}

// {
//   "vesta": 8,
//   "jupiter": 50,
//   "saturnus": 27,
//   "venus": 20,
//   "mercurius": 30,
//   "mars": 56,
//   "minerva": 15
// }

func ConcordiaScore(stats models.Stats) float64 {
	keys := []string{"vesta", "jupiter", "saturnus", "venus", "mercurius", "mars", "minerva"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	if val, ok := stats.Data["concordia"].(float64); ok {
		score += val
	}

	return score
}

func ConcordiaSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := ConcordiaScore(stats[i])
		score2 := ConcordiaScore(stats[j])
		return score1 < score2
	}
}

// {"7x7":true,"buttons":31,"empty":9}

func PatchworkScore(stats models.Stats) float64 {
	score := 0.0

	if val, ok := stats.Data["7x7"].(bool); ok {
		if val {
			score += 7
		}
	}

	score += stats.Data["buttons"].(float64)
	score -= 2 * stats.Data["empty"].(float64)

	return score
}

func PatchworkSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := PatchworkScore(stats[i])
		score2 := PatchworkScore(stats[j])
		return score1 < score2
	}
}

// {"armor":6,"loot":0,"points":99,"valkyrie":15}

func RaidersScore(stats models.Stats) float64 {
	keys := []string{"armor", "loot", "points", "valkyrie"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	return score
}

func RaidersSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := RaidersScore(stats[i])
		score2 := RaidersScore(stats[j])
		return score1 < score2
	}
}

// {"buildings":15,"cathedral":20,"debt":0,"gold-marble":2,"money":0,"prison":0,"virtue":7}

func ArchitectsScore(stats models.Stats) float64 {
	keys := []string{"buildings", "cathedral", "debt", "gold-marble", "money", "prison", "virtue"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	return score
}

func ArchitectsSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := ArchitectsScore(stats[i])
		score2 := ArchitectsScore(stats[j])

		if int64(score1) == int64(score2) {
			virtue1 := stats[i].Data["virtue"].(float64)
			virtue2 := stats[j].Data["virtue"].(float64)

			if int64(virtue1) == int64(virtue2) {
				money1 := stats[i].Data["money"].(float64)
				money2 := stats[j].Data["money"].(float64)

				return money1 < money2
			}

			return virtue1 < virtue2
		}

		return score1 < score2
	}
}

func KemetScore(stats models.Stats) float64 {
	keys := []string{"temple_temporary", "temple_permanent", "temple_allgods", "battle", "tile", "pyramid", "sanctuary", "sphinx"}

	score := 0.0
	for _, key := range keys {
		if val, ok := stats.Data[key].(float64); ok {
			score += val
		}
	}

	return score
}

func KemetSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := KemetScore(stats[i])
		score2 := KemetScore(stats[j])

		if int64(score1) == int64(score2) {
			battle1 := stats[i].Data["battle"].(float64)
			battle2 := stats[j].Data["battle"].(float64)

			if int64(battle1) == int64(battle2) {
				order1 := stats[i].Data["order"].(float64)
				order2 := stats[j].Data["order"].(float64)

				return order1 > order2
			}

			return battle1 < battle2
		}

		return score1 < score2
	}
}
