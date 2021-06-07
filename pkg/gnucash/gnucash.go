package gnucash

import (
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"time"
)

func countSince(createdAtTime time.Time) (int, int) {
	now := time.Now()
	months := 0
	days := 0
	month := createdAtTime.Month()
	day := createdAtTime.Day()

	for createdAtTime.Before(now) {
		createdAtTime = createdAtTime.Add(time.Hour * 24)

		nextMonth := createdAtTime.Month()
		if nextMonth != month {
			months++
		}

		nextDay := createdAtTime.Day()
		if nextDay != day {
			days++
		}

		month = nextMonth
		day = nextDay
	}

	return months, days
}

func GetExpenseByMonth(c *gin.Context) {
	expense := c.Param("expense")

	database, err := db.DatabaseConnect("gnucash")
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	rs, err := getExpenseByMonth(database, expense)
	if err != nil {
		util.Error(c, err)
		return
	}

	months, days := countSince(rs[0].Date)

	count := 0.0
	for _, item := range rs {
		count += item.Sum
	}

	util.Success(c, &map[string]interface{}{
		"data": rs,
		"calc": map[string]interface{}{
			"months":   months,
			"days":     days,
			"total":    count,
			"permonth": count / float64(months),
			"perday":   count / float64(days),
		},
	})
}
