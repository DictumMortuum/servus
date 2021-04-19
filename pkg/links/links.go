package links

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
)

func AddLink(c *gin.Context) {
	var form LinkRow

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

	exists, err := LinkExists(database, form)
	if err != nil {
		util.Error(c, err)
		return
	}

	payload["data"] = form
	payload["exists"] = exists

	if !exists {
		err = CreateLink(database, form)
		if err != nil {
			util.Error(c, err)
			return
		}
	}

	util.Success(c, &payload)
}
