package search

import (
	"github.com/gocolly/colly/v2"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	price = regexp.MustCompile("([0-9]+.[0-9]+)")
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

func Hamming(s1 string, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)
	distance := 0

	if len(s1) > len(s2) {
		r1 = []rune(s1)[:len(r2)]
	} else if len(s1) < len(s2) {
		r2 = []rune(s2)[:len(r1)]
	}

	if len(r1) != len(r2) {
		log.Println(s1)
		log.Println(s2)
		return -1
	}

	for i, v := range r1 {
		if r2[i] != v {
			distance += 1
		}
	}

	return distance
}
