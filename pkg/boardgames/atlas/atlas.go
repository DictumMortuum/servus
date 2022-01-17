package atlas

import (
	"encoding/json"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Search(name string) ([]models.AtlasResult, error) {
	var rs struct {
		Games []models.AtlasResult `json:"games"`
	}

	link := "https://api.boardgameatlas.com/api/search?name=" + url.QueryEscape(name) + "&pretty=true&client_id=" + config.App.Atlas.ClientId
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

func AtlasSearch2(game models.Boardgame) ([]models.AtlasResult, error) {
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

func AtlasSearch(c *gin.Context) {
	var rs struct {
		Games []models.AtlasResult `json:"games"`
	}

	type params struct {
		Term string `json:"term"`
	}

	var payload params
	c.BindJSON(&payload)

	link := "https://api.boardgameatlas.com/api/search?name=" + url.QueryEscape(payload.Term) + "&pretty=true&client_id=" + config.App.Atlas.ClientId
	fmt.Println(link)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		util.Error(c, err)
		return
	}

	conn := &http.Client{}

	resp, err := conn.Do(req)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error(c, err)
		return
	}

	err = json.Unmarshal(body, &rs)
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}
