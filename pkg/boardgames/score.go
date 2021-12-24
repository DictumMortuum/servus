package boardgames

import (
	"github.com/DictumMortuum/servus/pkg/models"
)

type column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Tiebreak string `json:"tiebreak,omitempty"`
}

type settings struct {
	Columns     []column `json:"columns,omitempty"`
	Cooperative string   `json:"cooperative,omitempty"`
}

func DatabaseScore(play models.Play) (func(models.Stats) float64, error) {
	var s settings

	err := play.BoardgameSettings.Unmarshal(&s)
	if err != nil {
		return nil, err
	}

	return func(stats models.Stats) float64 {
		score := 0.0
		base := 0.01

		for _, col := range s.Columns {
			switch col.Type {
			case "int":
				if val, ok := stats.Data[col.Name].(float64); ok {
					score += val
				}
			default:
				if val, ok := stats.Data[col.Name].(float64); ok {
					score += val
				}
			}

			if col.Tiebreak != "" {
				if val, ok := stats.Data[col.Name].(float64); ok {
					if col.Tiebreak == "asc" {
						score += val * base
					} else if col.Tiebreak == "desc" {
						score -= val * base
					}
				}
				base *= 0.1
			}
		}

		return score
	}, nil
}

func isCooperative(play models.Play) bool {
	if val, ok := play.BoardgameSettings["cooperative"]; ok {
		return val.(bool)
	}

	return false
}

func getFuncs(play models.Play) (func(models.Stats) float64, func([]models.Stats) func(i, j int) bool) {
	// here the boardgame has the settings
	if len(play.BoardgameSettings) != 0 {
		f, _ := DatabaseScore(play)
		return f, DefaultSort
	}

	switch play.BoardgameId {
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
	case 199792:
		return DefaultScore, DefaultSort
	case 169786:
		return DefaultScore, DefaultSort
	case 296151:
		return ViscountsScore, ViscountsSort
	case 164928:
		return OrleansScore, OrleansSort
	case 198994:
		return HealthScore, HealthSort
	case 193738:
		return GreatWesternTrailScore, DefaultSort
	case 301880:
		return DefaultScore, DefaultSort
	case 216132:
		return CaledoniaScore, DefaultSort
	case 276025:
		return DefaultScore, DefaultSort
	case 244521:
		return DefaultScore, DefaultSort
	case 2651:
		return DefaultScore, DefaultSort
	case 3076:
		return DefaultScore, DefaultSort
	case 275010:
		return DefaultScore, DefaultSort
	case 158899:
		return DefaultScore, DefaultSort
	case 253499:
		return WarOfWhispersScore, WarOfWhispersSort
	case 110277:
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

func HealthScore(stats models.Stats) float64 {
	return stats.Data["health"].(float64)
}

func HealthSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		return HealthScore(stats[i]) < HealthScore(stats[j])
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

	// tiebreak
	score += stats.Data["money"].(float64) * 0.01

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

//
// { "score": 5, "swaps?": 1 }
//

func WarOfWhispersScore(stats models.Stats) float64 {
	keys := []string{"score"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	if val, ok := stats.Data["swaps"]; ok {
		score -= val.(float64) * 0.01
	}

	return score
}

func WarOfWhispersSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := WarOfWhispersScore(stats[i])
		score2 := WarOfWhispersScore(stats[j])

		if int64(score1) == int64(score2) {
			swap1, ok1 := stats[i].Data["swaps"].(float64)
			swap2, ok2 := stats[j].Data["swaps"].(float64)

			if !ok1 && !ok2 {
				return true
			}

			if !ok1 && ok2 {
				return false
			}

			if ok1 && !ok2 {
				return true
			}

			return swap1 <= swap2
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

	score += stats.Data["coin_points"].(float64) * 0.01

	return score
}

func DuelSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		if val, ok := stats[i].Data["battle_victory"].(bool); ok {
			if val {
				return false
			}
		}

		if val, ok := stats[j].Data["battle_victory"].(bool); ok {
			if val {
				return true
			}
		}

		if val, ok := stats[i].Data["science_victory"].(bool); ok {
			if val {
				return false
			}
		}

		if val, ok := stats[j].Data["science_victory"].(bool); ok {
			if val {
				return true
			}
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
	keys := []string{"birds", "bonus", "endofround", "eggs", "food", "tucked", "nectar"}

	score := 0.0
	for _, key := range keys {
		if val, ok := stats.Data[key].(float64); ok {
			score += val
		}
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
		if val, ok := stats.Data[key].(float64); ok {
			score += val
		}
	}

	score += stats.Data["order"].(float64) * 0.01

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
			rows1, ok1 := stats[i].Data["rows"].(float64)
			rows2, ok2 := stats[j].Data["rows"].(float64)

			if ok1 && ok2 {
				return rows1 < rows2
			} else {
				return true
			}
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

	score += stats.Data["research"].(float64) * 0.01

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

// {"armor":6,"loot":5,"points":59,"valkyrie":7,"cards": 2}

func RaidersScore(stats models.Stats) float64 {
	keys := []string{"armor", "loot", "points", "valkyrie", "cards"}

	score := 0.0
	for _, key := range keys {
		if val, ok := stats.Data[key].(float64); ok {
			score += val
		}
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
	keys := []string{"buildings", "cathedral", "debt", "gold-marble", "money", "prison", "virtue", "bottom"}

	score := 0.0
	for _, key := range keys {
		if val, ok := stats.Data[key].(float64); ok {
			score += val
		}
	}

	score += stats.Data["virtue"].(float64) * 0.01

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

// {"build": 19, "castle": 28, "transcribe": 0, "castle_leader": 5, "deeds": 10, "dept": -2, "poverty": 12 }

func ViscountsScore(stats models.Stats) float64 {
	keys := []string{"build", "castle", "transcribe", "castle_leader", "deeds", "dept", "poverty"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	return score
}

func ViscountsSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := ViscountsScore(stats[i])
		score2 := ViscountsScore(stats[j])
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

	if val, ok := stats.Data["battle"]; ok {
		score += val.(float64) * 0.01
	}

	if val, ok := stats.Data["order"]; ok {
		score -= val.(float64) * 0.001
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

// {"grain": 3, "cheese": 6, "wine": 0, "wool": 12, "brocade": 10, "coin": 37, "citizen": 24, "buildings": 24}

func OrleansScore(stats models.Stats) float64 {
	keys := []string{"grain", "cheese", "wine", "wool", "brocade", "coin", "citizen", "buildings", "orders", "tiles"}

	score := 0.0
	for _, key := range keys {
		if val, ok := stats.Data[key].(float64); ok {
			score += val
		}
	}

	if val, ok := stats.Data["track"]; ok {
		score -= val.(float64) * 0.01
	}

	return score
}

func OrleansSort(stats []models.Stats) func(i, j int) bool {
	return func(i, j int) bool {
		score1 := OrleansScore(stats[i])
		score2 := OrleansScore(stats[j])

		if int64(score1) == int64(score2) {

			track1 := 99.0
			if val, ok := stats[i].Data["track"].(float64); ok {
				track1 = val
			}

			track2 := 99.0
			if val, ok := stats[j].Data["track"].(float64); ok {
				track2 = val
			}

			return track1 > track2
		}

		return score1 < score2
	}
}

func GreatWesternTrailScore(stats models.Stats) float64 {
	keys := []string{"dollars", "buildings", "cities", "stations", "hazards", "cards", "objectives", "tasks", "workers", "disk", "market"}

	score := 0.0
	for _, key := range keys {
		score += stats.Data[key].(float64)
	}

	return score
}

func CaledoniaScore(stats models.Stats) float64 {
	keys := []string{"points", "base", "manufactured", "coin", "grape", "cotton", "bananas", "bamboo", "exports", "settlement"}

	score := 0.0
	for _, key := range keys {
		if val, ok := stats.Data[key].(float64); ok {
			score += val
		}
	}

	return score
}
