package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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

func Success2(c *gin.Context, payload interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"error":    nil,
		"response": payload,
	})
}

func Success(c *gin.Context, payload *map[string]interface{}) {
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
