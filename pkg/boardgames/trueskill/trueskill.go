package trueskill

import (
	"github.com/DictumMortuum/servus/pkg/models"
	trueskill "github.com/mafredri/go-trueskill"
)

func Calculate(plays []models.Play) []models.Play {
	ts := trueskill.New()
	players := map[string]trueskill.Player{}

	for idx := range plays {
		if plays[idx].IsCooperative() {
			continue
		}

		// Reverse the array, so that winner is on the top.
		for i, j := 0, len(plays[idx].Stats)-1; i < j; i, j = i+1, j-1 {
			plays[idx].Stats[i], plays[idx].Stats[j] = plays[idx].Stats[j], plays[idx].Stats[i]
		}

		// Generate a list of the names of the players
		playersInPlay := []string{}
		for _, stat := range plays[idx].Stats {
			playersInPlay = append(playersInPlay, stat.Player)
		}

		// Generate a list of draws, if there are any
		draws := []bool{}
		flag := false
		for i := 0; i < len(plays[idx].Stats)-1; i++ {
			if plays[idx].Stats[i].Data["score"] == plays[idx].Stats[i+1].Data["score"] {
				draws = append(draws, true)
				flag = true
			} else {
				draws = append(draws, false)
			}
		}

		if flag {
			// fmt.Println(plays[idx])
			// fmt.Println(draws)
		}

		// For each player that participated, generate a trueskill structure - initialize it if it's his first time.
		var playerSkills []trueskill.Player
		for _, name := range playersInPlay {
			if val, ok := players[name]; ok {
				playerSkills = append(playerSkills, trueskill.NewPlayer(val.Mu(), val.Sigma()))
			} else {
				players[name] = ts.NewPlayer()
				playerSkills = append(playerSkills, players[name])
			}
		}

		adjustedPlayers, probability := ts.AdjustSkillsWithDraws(playerSkills, draws)
		// adjustedPlayers and player names are in order - useful to copy over the new ratings
		for i, name := range playersInPlay {
			previous_score := ts.TrueSkill(players[name])
			players[name] = trueskill.NewPlayer(adjustedPlayers[i].Mu(), adjustedPlayers[i].Sigma())
			plays[idx].Stats[i].Mu = adjustedPlayers[i].Mu()
			plays[idx].Stats[i].Sigma = adjustedPlayers[i].Sigma()
			plays[idx].Stats[i].TrueSkill = ts.TrueSkill(players[name])
			plays[idx].Stats[i].Delta = plays[idx].Stats[i].TrueSkill - previous_score
		}
		plays[idx].Probability = probability * 100
		plays[idx].Draws = draws
	}

	return plays
}
