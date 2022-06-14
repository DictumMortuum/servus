package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"time"
)

const (
	DAY_OFF         = -1
	NORMAL_LEAVE    = -2
	SICK_LEAVE      = -3
	PUBLIC_HOLIDAY  = -4
	PERENNIAL_LEAVE = -5
	MARRIAGE_LEAVE  = -6
	UNKNOWN_LEAVE   = -100
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
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
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

func FormatShiftColor(shift int) string {
	if shift == 3 {
		return "table-danger"
	} else if shift == 23 {
		return "table-primary"
	} else if shift == DAY_OFF {
		return "table-success"
	} else if shift == NORMAL_LEAVE {
		return "table-warning"
	} else if shift == SICK_LEAVE {
		return "table-warning"
	} else if shift == PUBLIC_HOLIDAY {
		return "table-warning"
	} else if shift == PERENNIAL_LEAVE {
		return "table-warning"
	} else if shift == MARRIAGE_LEAVE {
		return "table-warning"
	} else {
		return ""
	}
}

func FormatShift(shift int) string {
	if shift >= 0 {
		return fmt.Sprintf("Βάρδια %d", shift)
	} else if shift == DAY_OFF {
		return "Ρεπό"
	} else if shift == NORMAL_LEAVE {
		return "Άδεια"
	} else if shift == SICK_LEAVE {
		return "Ασθένεια"
	} else if shift == PUBLIC_HOLIDAY {
		return "Αργία"
	} else if shift == PERENNIAL_LEAVE {
		return "Πολυετία"
	} else if shift == MARRIAGE_LEAVE {
		return "Άδεια γάμου"
	} else {
		return "Άγνωστο"
	}
}
