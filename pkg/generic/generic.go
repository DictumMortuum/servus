package generic

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func A(f []func(*sqlx.DB, *models.QueryBuilder) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		args, err := models.NewArgsFromContext(c)
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

		rs := []interface{}{}

		for _, fn := range f {
			data, err := fn(database, args)
			if err != nil {
				util.Error(c, err)
				return
			}
			rs = append(rs, data)
		}

		c.JSON(http.StatusOK, rs)
	}
}

func F(f func(*sqlx.DB, *models.QueryBuilder) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		args, err := models.NewArgsFromContext(c)
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

		data, err := f(database, args)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

func S(d string, f func(*sqlx.DB, *models.QueryBuilder) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		args, err := models.NewArgsFromContext(c)
		if err != nil {
			util.Error(c, err)
			return
		}

		database, err := db.DatabaseConnect(d)
		if err != nil {
			util.Error(c, err)
			return
		}
		defer database.Close()

		data, err := f(database, args)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.JSON(http.StatusOK, data)
	}
}
