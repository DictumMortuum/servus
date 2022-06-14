package generic

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"sync"
)

func C(f []func(*sqlx.DB, *models.QueryBuilder) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		var wg sync.WaitGroup

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

		for i := 1; i <= len(f); i++ {
			wg.Add(1)
			fn := f[i-1]
			go func() {
				defer wg.Done()
				fn(database, args)
			}()
		}

		wg.Wait()

		c.JSON(http.StatusOK, nil)
	}
}

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
