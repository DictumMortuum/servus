package boardgames

import (
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/util"
	"github.com/gin-gonic/gin"
)

func GetDuel(c *gin.Context) {
	database, err := db.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	duels, err := getDuels(database)
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &duels)
}
