package links

import (
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/util"
	"github.com/gin-gonic/gin"
)

func AddLink(c *gin.Context) {
	var form db.LinkRow

	payload := map[string]interface{}{}

	err := c.ShouldBind(&form)
	if err != nil {
		util.Error(c, err)
		return
	}

	database, err := db.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	exists, err := db.LinkExists(database, form)
	if err != nil {
		util.Error(c, err)
		return
	}

	payload["data"] = form
	payload["exists"] = exists

	if !exists {
		err = db.CreateLink(database, form)
		if err != nil {
			util.Error(c, err)
			return
		}
	}

	util.Success(c, &payload)
}
