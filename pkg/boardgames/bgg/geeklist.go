package bgg

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GeeklistItem struct {
	XMLName  xml.Name `xml:"item" json:"-"`
	ObjectId int64    `xml:"objectid,attr" json:"id"`
	ItemId   int      `xml:"id,attr" json:"item_id"`
	Name     string   `xml:"objectname,attr" json:"name"`
	Date     string   `xml:"editdate,attr" json:"date"`
	Body     string   `xml:"body" json:"body"`
}

type GeeklistRs struct {
	XMLName xml.Name       `xml:"geeklist" json:"-"`
	Items   []GeeklistItem `xml:"item" json:"items"`
}

func Geeklist(id int64) (*GeeklistRs, error) {
	link := fmt.Sprintf("https://boardgamegeek.com/xmlapi2/geeklist/%d", id)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	conn := &http.Client{}

	resp, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tmp := GeeklistRs{}

	err = xml.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}

	return &tmp, nil
}
