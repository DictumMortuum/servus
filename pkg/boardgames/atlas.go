package boardgames

import (
	"encoding/json"
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/models"
	"io/ioutil"
	"net/http"
	"net/url"
)

func AtlasSearch(game models.Boardgame) ([]models.AtlasResult, error) {
	var rs struct {
		Games []models.AtlasResult `json:"games"`
	}

	link := "https://api.boardgameatlas.com/api/search?name=" + url.QueryEscape(game.Name) + "&pretty=true&client_id=" + config.App.Atlas.ClientId
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	conn := &http.Client{}

	resp, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &rs)
	if err != nil {
		return nil, err
	}

	return rs.Games, nil
}
