package generic

import (
	"fmt"
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

func Register(
	router *gin.RouterGroup,
	get func(*sqlx.DB, int64) (interface{}, error),
	getlist func(*sqlx.DB, Args) (interface{}, int, error),
	create func(*sqlx.DB, map[string]interface{}) (interface{}, error),
	update func(*sqlx.DB, int64, map[string]interface{}) (interface{}, error),
	delete func(*sqlx.DB, int64) (interface{}, error),
) {
	router.GET("/:id", _get(get))
	router.GET("", _getlist(getlist))
	router.POST("", _create(create))
	router.PUT("/:id", _update(update))
	router.DELETE("/:id", _delete(delete))
}

func _get(f func(*sqlx.DB, int64) (interface{}, error)) func(*gin.Context) {
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

func _getlist(f func(*sqlx.DB, Args) (interface{}, int, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		args := ParseArgs(c)

		database, err := db.Conn()
		if err != nil {
			util.Error(c, err)
			return
		}
		defer database.Close()

		data, length, err := f(database, args)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.Header("X-Total-Count", fmt.Sprintf("%d", length))
		c.JSON(http.StatusOK, gin.H{
			"data":  data,
			"total": length,
		})
	}
}

func _create(f func(*sqlx.DB, map[string]interface{}) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		var args map[string]interface{}
		c.BindJSON(&args)

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

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func _update(f func(*sqlx.DB, int64, map[string]interface{}) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		var args map[string]interface{}
		c.BindJSON(&args)
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

		data, err := f(database, id, args)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func _delete(f func(*sqlx.DB, int64) (interface{}, error)) func(*gin.Context) {
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
