package boardgames

import (
	"github.com/DictumMortuum/servus/pkg/models"
)

type column struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Tiebreak string  `json:"tiebreak,omitempty"`
	Factor   float64 `json:"factor,omitempty"`
	Exclude  bool    `json:"exclude,omitempty"`
}

type settings struct {
	Columns     []column `json:"columns,omitempty"`
	Cooperative string   `json:"cooperative,omitempty"`
	Autowin     []string `json:"autowin,omitempty"`
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
			if col.Tiebreak != "" {
				factor := col.Factor
				if factor == 0 {
					factor = base
				}

				if val, ok := stats.Data[col.Name].(float64); ok {
					if col.Tiebreak == "asc" {
						score += val * base
					} else if col.Tiebreak == "desc" {
						score -= val * base
					}
				}
				base *= 0.1
			}

			if col.Exclude {
				continue
			}

			switch col.Type {
			case "int":
				if val, ok := stats.Data[col.Name].(float64); ok {
					score += val
				}
			case "negint":
				if val, ok := stats.Data[col.Name].(float64); ok {
					score -= val
				}
			default:
				if val, ok := stats.Data[col.Name].(float64); ok {
					score += val
				}
			}
		}

		return score
	}, nil
}

func DatabaseSort(play models.Play) (func([]models.Stats) func(i, j int) bool, error) {
	var s settings

	err := play.BoardgameSettings.Unmarshal(&s)
	if err != nil {
		return nil, err
	}

	return func(stats []models.Stats) func(i, j int) bool {
		return func(i, j int) bool {
			for _, key := range s.Autowin {
				val1, ok1 := stats[i].Data[key].(bool)
				val2, ok2 := stats[j].Data[key].(bool)

				if !ok1 && !ok2 {
					continue
				}

				if ok1 && !ok2 {
					if val1 {
						return false
					}
				}

				if !ok1 && ok2 {
					if val2 {
						return true
					}
				}

				if ok1 && ok2 {
					if val1 && val2 {
						continue
					}

					if val1 {
						return false
					}

					if val2 {
						return true
					}
				}
			}

			return DefaultScore(stats[i]) < DefaultScore(stats[j])
		}
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
		g, _ := DatabaseSort(play)
		return f, g
	}

	return nil, nil
}

func DefaultScore(stats models.Stats) float64 {
	return stats.Data["score"].(float64)
}
