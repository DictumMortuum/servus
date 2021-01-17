package weight

import (
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/util"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func AddWeight(c *gin.Context) {
	var data WeightRow

	err := c.ShouldBind(&data)
	if err != nil {
		util.Error(c, err)
		return
	}

	data.Date, err = time.Parse("2006-01-02 15:04:05", data.DateRaw)
	if err != nil {
		util.Error(c, err)
		return
	}
	log.Println(data)
	database, err := db.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	err = CreateWeight(database, data)
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &data)
}
