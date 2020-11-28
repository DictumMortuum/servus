package zerotier

import (
	"bytes"
	"encoding/json"
	"github.com/DictumMortuum/servus/config"
	"github.com/DictumMortuum/servus/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	endpoint = "https://my.zerotier.com/api"
)

func GetNode(c *gin.Context) {
	network := config.App.Zerotier.Network
	member := c.Param("member")

	req, err := http.NewRequest("GET", endpoint+"/network/"+network+"/member/"+member, nil)
	if err != nil {
		util.Error(c, err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+config.App.Zerotier.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
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

	retval := map[string]interface{}{}
	err = json.Unmarshal(body, &retval)
	if err != nil {
		util.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"error":    nil,
		"response": retval,
	})
}

func PostNode(c *gin.Context) {
	network := config.App.Zerotier.Network
	member := c.PostForm("member")

	payload := map[string]interface{}{
		"config": map[string]interface{}{
			"authorized": true,
		},
		"description": "joined " + time.Now().Format("02-01-2006"),
	}

	jsonStr, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint+"/network/"+network+"/member/"+member, bytes.NewBuffer(jsonStr))
	if err != nil {
		util.Error(c, err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+config.App.Zerotier.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
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

	retval := map[string]interface{}{}
	err = json.Unmarshal(body, &retval)
	if err != nil {
		util.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"error":    nil,
		"response": retval,
	})
}
