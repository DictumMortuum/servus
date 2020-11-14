package gas

import (
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Render(c *gin.Context) {
	database, err := db.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	gas, err := db.GetGas(database)
	if err != nil {
		util.Error(c, err)
		return
	}

	state := gin.H{
		"title": "Gas",
		"primary": map[string]interface{}{
			"enabled": true,
			"desc":    "Add",
		},
		"secondary": map[string]interface{}{
			"enabled": false,
		},
		"gas":   gas,
		"error": nil,
	}

	c.HTML(http.StatusOK, "gas.html", state)
}

func AddFuelStats(c *gin.Context) {
	var form db.FuelStatsRow

	err := c.ShouldBind(&form)
	if err != nil {
		util.Error(c, err)
		return
	}

	years := map[string]interface{}{
		"a": form.Kilometers,
		"b": form.LitreAverage,
		"c": form.Duration,
		"d": form.KilometersPerHour,
	}

	util.Success(c, &years)
}
