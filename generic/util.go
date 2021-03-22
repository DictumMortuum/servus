package generic

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type Args struct {
	Sort      string
	Order     string
	Range     []int64
	Page      int64
	RefKey    string
	Id        []int64
	FilterKey string
	FilterVal string
}

func ParseArgs(c *gin.Context) Args {
	args := c.Request.URL.Query()
	rs := Args{
		Page: 5,
	}

	if val, ok := args["_sort"]; ok {
		rs.Sort = val[0]
	}

	if val, ok := args["_order"]; ok {
		rs.Order = val[0]
	} else {
		rs.Order = "ASC"
	}

	if val, ok := args["_start"]; ok {
		start, err := strconv.ParseInt(val[0], 10, 64)
		if err == nil {
			rs.Range = append(rs.Range, start)
		}
	}

	if val, ok := args["_end"]; ok {
		end, err := strconv.ParseInt(val[0], 10, 64)
		if err == nil {
			rs.Range = append(rs.Range, end)
		}
	}

	if len(rs.Range) == 2 {
		rs.Page = rs.Range[1] - rs.Range[0] + 1
	}

	if val, ok := args["id"]; ok {
		for _, raw := range val {
			id, err := strconv.ParseInt(raw, 10, 64)
			if err == nil {
				rs.Id = append(rs.Id, id)
			}
			rs.RefKey = "id"
		}
	}

	for key, val := range args {
		if strings.HasSuffix(key, "id") && key != "id" {
			// this is a reference
			id, err := strconv.ParseInt(val[0], 10, 64)
			if err == nil {
				rs.RefKey = key
				rs.Id = append(rs.Id, id)
			}
		} else if !strings.HasPrefix(key, "_") && key != "id" {
			// this is a filter
			rs.FilterKey = key
			rs.FilterVal = val[0]
		}
	}

	return rs
}
