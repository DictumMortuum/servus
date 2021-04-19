package gnucash

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
)

func GetExpenseByMonth(c *gin.Context) {
	expense := c.Param("expense")

	database, err := db.DatabaseConnect("gnucash")
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	rs, err := getExpenseByMonth(database, expense)
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}
