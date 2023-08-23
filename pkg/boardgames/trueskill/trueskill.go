package trueskill

import (
	"sort"

	"github.com/DictumMortuum/servus/pkg/models"
	trueskill "github.com/mafredri/go-trueskill"
)

func calculateWithoutTeams(ts *trueskill.Config, play *models.Play, players map[int64]trueskill.Player) *models.Play {
	// Reverse the array, so that winner is on the top.
	for i, j := 0, len(play.Stats)-1; i < j; i, j = i+1, j-1 {
		play.Stats[i], play.Stats[j] = play.Stats[j], play.Stats[i]
	}

	// Generate a list of the names of the players
	playersInPlay := []int64{}
	for _, stat := range play.Stats {
		playersInPlay = append(playersInPlay, stat.PlayerId)
	}

	// Generate a list of draws, if there are any
	draws := []bool{}
	flag := false
	for i := 0; i < len(play.Stats)-1; i++ {
		if play.Stats[i].Data["score"] == play.Stats[i+1].Data["score"] {
			draws = append(draws, true)
			flag = true
		} else {
			draws = append(draws, false)
		}
	}

	if flag {
		// fmt.Println(play)
		// fmt.Println(draws)
	}

	// For each player that participated, generate a trueskill structure - initialize it if it's his first time.
	var playerSkills []trueskill.Player
	for _, id := range playersInPlay {
		if val, ok := players[id]; ok {
			playerSkills = append(playerSkills, trueskill.NewPlayer(val.Mu(), val.Sigma()))
		} else {
			players[id] = ts.NewPlayer()
			playerSkills = append(playerSkills, players[id])
		}
	}

	adjustedPlayers, probability := ts.AdjustSkillsWithDraws(playerSkills, draws)
	// adjustedPlayers and player names are in order - useful to copy over the new ratings
	for i, id := range playersInPlay {
		previous_score := ts.TrueSkill(players[id])
		players[id] = trueskill.NewPlayer(adjustedPlayers[i].Mu(), adjustedPlayers[i].Sigma())
		play.Stats[i].Mu = adjustedPlayers[i].Mu()
		play.Stats[i].Sigma = adjustedPlayers[i].Sigma()
		play.Stats[i].TrueSkill = ts.TrueSkill(players[id])
		play.Stats[i].Delta = play.Stats[i].TrueSkill - previous_score
	}

	play.Probability = probability * 100
	play.Draws = draws
	return play
}

func generateTeam(ts *trueskill.Config, players map[int64]trueskill.Player, team []int64) trueskill.Player {
	mu := 0.0
	sigma := 0.0

	for _, id := range team {
		_, ok := players[id]
		if !ok {
			players[id] = ts.NewPlayer()
		}
	}

	for _, id := range team {
		player := players[id]
		mu += player.Mu()
		sigma += player.Sigma()
	}

	mu /= float64(len(team))
	sigma /= float64(len(team))

	return trueskill.NewPlayer(mu, sigma)
}

func getTeamScore(play models.Play, team []int64) float64 {
	first := team[0]

	for _, stat := range play.Stats {
		if stat.PlayerId == first {
			return stat.Data["score"].(float64)
		}
	}

	return -1
}

func calculateWithTeams(ts *trueskill.Config, play *models.Play, players map[int64]trueskill.Player) *models.Play {
	m := map[int64]int{}
	for i, stat := range play.Stats {
		m[stat.PlayerId] = i
	}

	type Data struct {
		Trueskill  trueskill.Player
		Adjustment trueskill.Player
		Score      float64
		Team       []int64
	}

	data := []Data{}
	for _, team := range play.GetTeams("teams") {
		data = append(data, Data{
			Trueskill: generateTeam(ts, players, team),
			Score:     getTeamScore(*play, team),
			Team:      team,
		})
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Score > data[j].Score
	})

	newTeams := [][]int64{}
	for _, team := range data {
		newTeams = append(newTeams, team.Team)
	}

	// Generate a list of draws, if there are any
	draws := []bool{}
	for i := 0; i < len(data)-1; i++ {
		if data[i].Score == data[i+1].Score {
			draws = append(draws, true)
		} else {
			draws = append(draws, false)
		}
	}

	// For each player that participated, generate a trueskill structure - initialize it if it's his first time.
	var teamSkills []trueskill.Player
	for _, team := range data {
		teamSkills = append(teamSkills, team.Trueskill)
	}

	adjustedTeamSkills, probability := ts.AdjustSkillsWithDraws(teamSkills, draws)
	for i, adjustment := range adjustedTeamSkills {
		data[i].Adjustment = adjustment
		mu := (adjustment.Mu() - data[i].Trueskill.Mu()) / float64(len(data[i].Team))
		sigma := (adjustment.Sigma() - data[i].Trueskill.Sigma()) / float64(len(data[i].Team))

		for _, id := range data[i].Team {
			previous_score := ts.TrueSkill(players[id])
			players[id] = trueskill.NewPlayer(players[id].Mu()+mu, players[id].Sigma()+sigma)
			play.Stats[m[id]].Mu = players[id].Mu()
			play.Stats[m[id]].Sigma = players[id].Sigma()
			play.Stats[m[id]].TrueSkill = ts.TrueSkill(players[id])
			play.Stats[m[id]].Delta = play.Stats[m[id]].TrueSkill - previous_score
		}
	}

	play.Probability = probability * 100
	play.Draws = draws
	play.PlaySettings["teams2"] = newTeams
	return play
}

func Calculate(plays []models.Play) []models.Play {
	ts := trueskill.New()
	players := map[int64]trueskill.Player{}
	rs := []models.Play{}

	for idx := range plays {
		if plays[idx].IsCooperative() {
			continue
		}

		teams := plays[idx].GetTeams("teams")
		if len(teams) > 0 {
			rs = append(rs, *calculateWithTeams(&ts, &plays[idx], players))
		} else {
			rs = append(rs, *calculateWithoutTeams(&ts, &plays[idx], players))
		}
	}

	return rs
}
