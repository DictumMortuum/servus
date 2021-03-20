package generic

import (
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

func GetOne(f func(*sqlx.DB, int64) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		arg := c.Params.ByName("id")

		id, err := strconv.ParseInt(arg, 10, 64)
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

		data, err := f(database, id)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func GetMany(f func(*sqlx.DB, []int64) (interface{}, int, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		req := c.Request.URL.Query()
		args := req["id"]
		ids := []int64{}

		for _, raw := range args {
			id, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				util.Error(c, err)
				return
			}

			ids = append(ids, id)
		}

		database, err := db.Conn()
		if err != nil {
			util.Error(c, err)
			return
		}
		defer database.Close()

		data, length, err := f(database, ids)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  data,
			"total": length,
		})
	}
}
