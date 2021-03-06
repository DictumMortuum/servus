package calendar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	TEMPORARY       = "/tmp/cal"
	DAY_OFF         = -1
	NORMAL_LEAVE    = -2
	SICK_LEAVE      = -3
	PUBLIC_HOLIDAY  = -4
	PERENNIAL_LEAVE = -5
	MARRIAGE_LEAVE  = -6
	UNKNOWN_LEAVE   = -100
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

func FormatCoworkers(coworkers []string) string {
	l := len(coworkers)
	var col []string

	for _, coworker := range coworkers {
		col = append(col, strings.Title(strings.ToLower(strings.TrimSpace(coworker))))
	}

	if l == 0 {
		return ""
	} else {
		return strings.Join(col, ", ")
	}
}

func RawToAbsence(raw string) int {
	re1 := regexp.MustCompile("[ΡPΕEΠΟO]{4}")
	re2 := regexp.MustCompile("[ΚΑKA]{2}")
	re3 := regexp.MustCompile("[ΑAΣΘ]{3}")
	re4 := regexp.MustCompile("[KΚ]{1}")
	re5 := regexp.MustCompile("[ΠOΟΛ]{3}")
	re6 := regexp.MustCompile("[AΑΓ]{2}")

	if re1.MatchString(raw) {
		return DAY_OFF
	} else if re2.MatchString(raw) {
		return NORMAL_LEAVE
	} else if re3.MatchString(raw) {
		return SICK_LEAVE
	} else if re4.MatchString(raw) {
		return PUBLIC_HOLIDAY
	} else if re5.MatchString(raw) {
		return PERENNIAL_LEAVE
	} else if re6.MatchString(raw) {
		return MARRIAGE_LEAVE
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
