package generic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RawArgs struct {
	Sort   string `form:"sort" binding:"omitempty"`
	Range  string `form:"range" binding:"omitempty"`
	Filter string `form:"filter" binding:"omitempty"`
}

type Args struct {
	Sort   []string
	Range  []int
	Filter map[string]interface{}
	Page   int
}

func ParseArgs(c *gin.Context) (error, *Args) {
	var rawArgs RawArgs
	var args Args

	err := c.ShouldBindQuery(&rawArgs)
	if err != nil {
		return err, nil
	}

	if rawArgs.Sort != "" {
		err = json.Unmarshal([]byte(rawArgs.Sort), &args.Sort)
		if err != nil {
			return err, nil
		}
	}

	if rawArgs.Range != "" {
		err = json.Unmarshal([]byte(rawArgs.Range), &args.Range)
		if err != nil {
			return err, nil
		}

		args.Page = args.Range[1] - args.Range[0] + 1
	}

	if rawArgs.Filter != "" {
		err = json.Unmarshal([]byte(rawArgs.Filter), &args.Filter)
		if err != nil {
			return err, nil
		}
	}

	return nil, &args
}

func Fail(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":   http.StatusBadRequest,
		"message":  err.Error(),
		"response": nil,
	})
}
