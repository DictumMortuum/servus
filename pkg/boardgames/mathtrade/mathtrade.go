package mathtrade

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	re       = regexp.MustCompile(`\([A-Z0-9_ ]+\) ([0-9A-Z\-]+)[ ]+receives \([A-Z0-9_ ]+\) ([0-9A-Z\-]+)`)
	wants_re = regexp.MustCompile(`([0-9A-Z\-]+) ==> (.*)`)
)

func dfs(key string, data map[string]string, visited *map[string]bool, rs *[]string) {
	if (*visited)[key] {
		(*rs) = append(*rs, key)
	}

	neighbour, ok := data[key]
	if !ok {
		return
	}

	if (*visited)[neighbour] {
		return
	} else {
		(*visited)[neighbour] = true
		dfs(neighbour, data, visited, rs)
	}

	return
}

func ParseResults(raw string) [][]string {
	result := map[string]string{}
	lines := strings.Split(raw, "\n")

	for _, line := range lines {
		if strings.Contains(line, "receives") && !strings.Contains(line, "sends") {
			rs := re.FindStringSubmatch(line)

			if len(rs) == 3 {
				result[rs[1]] = rs[2]
			}
		}
	}

	visited := map[string]bool{}

	for key := range result {
		visited[key] = false
	}

	retval := [][]string{}
	for key := range result {
		if !visited[key] {
			chain := []string{}
			dfs(key, result, &visited, &chain)
			retval = append(retval, chain)
		}
	}

	return retval
}

func ParseWants(raw string) map[string]string {
	lines := strings.Split(raw, "\n")
	rs := map[string]string{}

	for _, line := range lines {
		if strings.Contains(line, "==>") {
			wants_rs := wants_re.FindStringSubmatch(line)

			if len(wants_rs) == 3 {
				rs[wants_rs[1]] = wants_rs[2]
			}
		}
	}

	return rs
}

func ResultsOfficial(id string) (string, error) {
	link := fmt.Sprintf("https://bgg.activityclub.org/olwlg/%s-results-official.txt", id)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}

	conn := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := conn.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func OfficialWants(id string) (string, error) {
	link := fmt.Sprintf("https://bgg.activityclub.org/olwlg/%s-officialwants.txt", id)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}

	conn := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := conn.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func CountEur(data []string) int {
	count := 0

	for _, item := range data {
		xlate := strings.ToLower(item)

		if strings.Contains(xlate, "eur") || strings.Contains(xlate, "ευρ") {
			count++
		}
	}

	return count
}

func Analyze(c *gin.Context) {
	id := c.Param("id")

	raw_wants, _ := OfficialWants(id)
	raw_results, _ := ResultsOfficial(id)

	wants := ParseWants(raw_wants)
	results := ParseResults(raw_results)

	type Chain struct {
		Id        string   `json:"id"`
		Length    int      `json:"length"`
		Items     []string `json:"items"`
		RealMoney int      `json:"cash"`
	}

	retval := []Chain{}

	for _, chain := range results {
		rs := []string{}

		for _, item := range chain {
			rs = append(rs, fmt.Sprintf("%s - %s", item, wants[item]))
		}

		retval = append(retval, Chain{
			Id:        id,
			Length:    len(rs),
			Items:     rs,
			RealMoney: CountEur(rs),
		})
	}

	c.JSON(http.StatusOK, retval)
}
