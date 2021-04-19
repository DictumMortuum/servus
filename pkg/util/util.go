package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"time"
)

func stringToDateTimeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
		return time.Parse(time.RFC3339, data.(string))
	}

	return data, nil
}

func Default(f func() ([]interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		rs, err := f()
		if err != nil {
			Error(c, err)
			return
		}

		Success(c, &rs)
	}
}

func Render(c *gin.Context, data gin.H, templateName string) {
	//switch c.Request.Header.Get("Accept") {
	//case "application/json":
	c.HTML(http.StatusOK, templateName, data)
}

func Error(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":   "BAD",
		"error":    err.Error(),
		"response": nil,
	})
}

func ErrorRedirect(c *gin.Context, endpoint string, err error) {
	c.SetCookie("error", err.Error(), 3600, "", "", false, true)
	c.Redirect(http.StatusMovedPermanently, endpoint)
}

func Success(c *gin.Context, payload interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"error":    nil,
		"response": payload,
	})
}

func FormatDate(t time.Time) string {
	return t.Format("02-01-2006")
}

func FormatDay(t time.Time) string {
	return t.Format("Monday")
}
