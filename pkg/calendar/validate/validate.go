package validate

import (
	"github.com/DictumMortuum/servus/pkg/calendar"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
)

func Weeks(c *calendar.Person) [][]int {
	weeks := [][]int{}
	retval := []int{}
	flag := false

	for i := 0; i < len(c.Calendar); i++ {
		if c.Calendar[i].DayName == "ΔΕ" && !flag {
			flag = true
			retval = append(retval, c.Calendar[i].Index)
		}

		if c.Calendar[i].DayName == "ΚΥ" && flag {
			flag = false
			retval = append(retval, c.Calendar[i].Index)
			weeks = append(weeks, retval)
			retval = []int{}
		}
	}

	return weeks
}

func DaysOff(c *calendar.Person) int {
	i := 0

	for _, day := range c.Calendar {
		if calendar.RawToShift(day.Raw) == calendar.DAY_OFF {
			i++
		}
	}

	return i
}

func DaysOffAfterNights(c *calendar.Person) bool {
	length := len(c.Calendar)

	for i := 0; i < length; i++ {
		if calendar.RawToShift(c.Calendar[i].Raw) == 23 {
			if i+1 < length {
				next := calendar.RawToShift(c.Calendar[i+1].Raw)

				if next == 23 || next == calendar.DAY_OFF {
					continue
				} else {
					return false
				}
			}
		}
	}

	return true
}

func MorningAfterEvening(c *calendar.Person) bool {
	length := len(c.Calendar)

	for i := 0; i < length; i++ {
		if calendar.RawToShift(c.Calendar[i].Raw) == 3 {
			if i+1 < length {
				next := calendar.RawToShift(c.Calendar[i+1].Raw)

				if next == 7 {
					return true
				}
			}
		}
	}

	return false
}

func DoubleDaysOff(c *calendar.Person) int {
	weeks := [][]int{}
	retval := []int{}
	flag := false

	for i := 0; i < len(c.Calendar); i++ {
		if calendar.RawToShift(c.Calendar[i].Raw) == calendar.DAY_OFF && !flag {
			flag = true
			retval = append(retval, c.Calendar[i].Index)
		}

		if calendar.RawToShift(c.Calendar[i].Raw) != calendar.DAY_OFF && flag {
			flag = false

			if retval[0]+1 < i {
				retval = append(retval, c.Calendar[i].Index-1)
				weeks = append(weeks, retval)
			}

			retval = []int{}
		}
	}

	return len(weeks)
}

func ValidateSheet(s *calendar.Sheet) map[string]interface{} {
	retval := map[string]interface{}{}

	for _, c := range s.People {
		r := map[string]interface{}{}

		r["name"] = c.Name
		r["surname"] = c.Surname
		r["id"] = c.Id
		r["speciality"] = c.Speciality

		checks := map[string]interface{}{}

		checks["daysoff"] = DaysOff(&c)
		checks["dayoffafternights"] = DaysOffAfterNights(&c)
		checks["contiguousdaysoffinmonth"] = DoubleDaysOff(&c)
		checks["morningshiftafterevenings"] = MorningAfterEvening(&c)

		r["checks"] = checks
		retval[c.Id] = r
	}

	return retval
}

func Validate(c *gin.Context) {
	file, _ := c.FormFile("files")
	c.SaveUploadedFile(file, calendar.TEMPORARY)

	sheet := calendar.Sheet{}
	err := sheet.New(calendar.TEMPORARY)
	data := ValidateSheet(&sheet)

	if err != nil {
		util.Error(c, err)
	} else {
		util.Success(c, &data)
	}
}
