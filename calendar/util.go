package calendar

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	TEMPORARY      = "/tmp/cal"
	DAY_OFF        = -1
	NORMAL_LEAVE   = -2
	SICK_LEAVE     = -3
	PUBLIC_HOLIDAY = -4
	UNKNOWN_LEAVE  = -100
)

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
	} else {
		return "Άγνωστο"
	}
}

func RawToAbsence(raw string) int {
	re1 := regexp.MustCompile("[ΡPΕEΠΟO]{4}")
	re2 := regexp.MustCompile("[ΚΑKA]{2}")
	re3 := regexp.MustCompile("[ΑAΣΘ]{3}")
	re4 := regexp.MustCompile("[KΚ]{1}")

	if re1.MatchString(raw) {
		return DAY_OFF
	} else if re2.MatchString(raw) {
		return NORMAL_LEAVE
	} else if re3.MatchString(raw) {
		return SICK_LEAVE
	} else if re4.MatchString(raw) {
		return PUBLIC_HOLIDAY
	} else {
		return UNKNOWN_LEAVE
	}
}

func RawToShift(raw string) int {
	re := regexp.MustCompile("[MHΜΗ]{2}([0-9][0-9])")
	s := re.FindStringSubmatch(raw)

	if s != nil {
		i, _ := strconv.Atoi(s[1])
		return i
	} else {
		return RawToAbsence(raw)
	}
}
