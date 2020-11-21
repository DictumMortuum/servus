package links

import (
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/util"
	"github.com/gin-gonic/gin"
)

func AddLink(c *gin.Context) {
	var form db.LinkRow

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

	err = db.CreateLink(database, form)
	if err != nil {
		util.Error(c, err)
		return
	}

	payload := map[string]interface{}{
		"data": form,
	}

	util.Success(c, &payload)
}
