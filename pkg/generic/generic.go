package generic

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GET(p models.Getable) func(*gin.Context) {
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

		data, err := p.Get(database, id)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

func GETLIST(p models.Getlistable) func(*gin.Context) {
	return func(c *gin.Context) {
		args := ParseArgs(c)

		database, err := db.Conn()
		if err != nil {
			util.Error(c, err)
			return
		}
		defer database.Close()

		rs, count, err := p.GetList(database, args)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.Header("X-Total-Count", fmt.Sprintf("%d", count))
		c.JSON(http.StatusOK, rs)
	}
}

func POST(p models.Createable) func(*gin.Context) {
	return func(c *gin.Context) {
		var args map[string]interface{}
		c.BindJSON(&args)

		database, err := db.Conn()
		if err != nil {
			util.Error(c, err)
			return
		}
		defer database.Close()

		qb := models.QueryBuilder{
			Columns: []string{},
		}

		for key := range args {
			qb.Columns = append(qb.Columns, key)
		}

		qb.Data = args

		data, err := p.Create(database, qb)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

func PUT(p models.Updateable) func(*gin.Context) {
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

		qb := models.QueryBuilder{
			Columns: []string{},
		}

		for key := range args {
			qb.Columns = append(qb.Columns, key)
		}

		qb.Data = args

		data, err := p.Update(database, id, qb)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

func DELETE(p models.Deleteable) func(*gin.Context) {
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

		data, err := p.Delete(database, id)
		if err != nil {
			util.Error(c, err)
			return
		}

		c.JSON(http.StatusOK, data)
	}
}
