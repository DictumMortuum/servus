package search

import (
	"github.com/gocolly/colly/v2"
	"mvdan.cc/xurls/v2"
	"regexp"
	"strconv"
	"strings"
)

var (
	price = regexp.MustCompile("([0-9]+[,.][0-9]+)")
)

func hasClass(e *colly.HTMLElement, c string) bool {
	raw := e.Attr("class")
	classes := strings.Split(raw, " ")

	for _, class := range classes {
		if class == c {
			return true
		}
	}

	return false
}

func childHasClass(e *colly.HTMLElement, child string, c string) bool {
	raw := e.ChildAttr(child, "class")
	classes := strings.Split(raw, " ")

	for _, class := range classes {
		if class == c {
			return true
		}
	}

	return false
}

func getPrice(raw string) float64 {
	raw = strings.ReplaceAll(raw, ",", ".")
	match := price.FindStringSubmatch(raw)

	if len(match) > 0 {
		price, _ := strconv.ParseFloat(match[1], 64)
		return price
	} else {
		return 0.0
	}
}

func getURL(raw string) []string {
	xurl := xurls.Strict()
	return xurl.FindAllString(raw, -1)
}

func Hamming(s1 string, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)
	var column []bool

	if len(r1) >= len(r2) {
		column = make([]bool, len(r1)+1)

		for i := 0; i < len(r2); i++ {
			column[i] = r1[i] == r2[i]
		}
	} else {
		column = make([]bool, len(r2)+1)

		for i := 0; i < len(r1); i++ {
			column[i] = r1[i] == r2[i]
		}
	}

	distance := 0

	for _, item := range column {
		if item == true {
			distance += 1
		}
	}

	return distance
}
