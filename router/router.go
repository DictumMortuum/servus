package router

import (
	"errors"
	"fmt"
	"github.com/DictumMortuum/servus/db"
	"github.com/DictumMortuum/servus/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}*/

type params struct {
	CheckUptime bool `form:"check_uptime"`
}

func Get(c *gin.Context) {
	var retval RouterRow
	var p params

	err := c.ShouldBind(&p)
	if err != nil {
		util.Error(c, err)
		return
	}

	req, err := http.NewRequest("GET", "http://192.168.1.1/comm/wan_cfg.sjs", nil)
	if err != nil {
		util.Error(c, err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "ID=dfuser")

	res1, err := http.DefaultClient.Do(req)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer res1.Body.Close()
	if res1.StatusCode != 200 {
		util.Error(c, errors.New(fmt.Sprintf("status code error: %d %s", res1.StatusCode, res1.Status)))
		return
	}

	bodyBytes, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		util.Error(c, err)
		return
	}

	vm := otto.New()

	_, err = vm.Run(string(bodyBytes))
	if err != nil {
		util.Error(c, err)
		return
	}

	data, err := vm.Run("PPP_ConnectionTable[0].TxBytes")
	if err != nil {
		util.Error(c, err)
		return
	}

	retval.DataUp, _ = data.ToInteger()

	data, err = vm.Run("PPP_ConnectionTable[0].RxBytes")
	if err != nil {
		util.Error(c, err)
		return
	}

	retval.DataDown, _ = data.ToInteger()

	data, err = vm.Run("PPP_ConnectionTable[0].UpTime")
	if err != nil {
		util.Error(c, err)
		return
	}

	uptime, _ := data.ToInteger()
	retval.Uptime = uptime
	t := time.Now().Add(time.Duration(-uptime) * time.Second)
	retval.Date = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())

	res2, err := http.Get("http://192.168.1.1/broadband/bd_dsl_detail.shtml?be=0&l0=2&l1=0&dtl=dt")
	if err != nil {
		util.Error(c, err)
		return
	}
	defer res2.Body.Close()
	if res2.StatusCode != 200 {
		util.Error(c, errors.New(fmt.Sprintf("status code error: %d %s", res2.StatusCode, res2.Status)))
		return
	}

	doc, err := goquery.NewDocumentFromReader(res2.Body)
	if err != nil {
		util.Error(c, err)
		return
	}

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_MAXBDWIDTH] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.MaxUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.MaxDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	})

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_BDWIDTH] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.CurrentUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.CurrentDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	})

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_CE] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.CRCUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.CRCDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	})

	doc.Find("td[key=PAGE_BD_DSL_DETAIL_FE] + td").Each(func(i int, s *goquery.Selection) {
		current := strings.Split(s.Text(), "/")
		retval.FECUp, _ = strconv.Atoi(strings.TrimSpace(current[0]))
		retval.FECDown, _ = strconv.Atoi(strings.TrimSpace(current[1]))
	})

	firstScript := doc.Find("script[language=javascript]").First()

	vm = otto.New()

	_, err = vm.Run("function GetWanDSLStatus(){}")
	if err != nil {
		util.Error(c, err)
		return
	}

	_, err = vm.Run(firstScript.Text())
	if err != nil {
		util.Error(c, err)
		return
	}

	snr, err := vm.Get("usNoiseMargin")
	if err != nil {
		util.Error(c, err)
		return
	}

	retval.SNRUp, _ = snr.ToInteger()

	snr, err = vm.Get("dsNoiseMargin")
	if err != nil {
		util.Error(c, err)
		return
	}

	retval.SNRDown, _ = snr.ToInteger()

	database, err := db.Conn()
	if err != nil {
		util.Error(c, err)
		return
	}
	defer database.Close()

	update := true

	if p.CheckUptime {
		exists, err := RouterExists(database, retval)
		if err != nil {
			util.Error(c, err)
			return
		}
		update = !exists
	}

	if update {
		err = CreateRouter(database, retval)
		if err != nil {
			util.Error(c, err)
			return
		}
	}

	util.Success2(c, &retval)
}
