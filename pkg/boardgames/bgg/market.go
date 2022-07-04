package bgg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Rs struct {
	Price string
}

func MarketProduct(id int64) (*Rs, error) {
	var conn = &http.Client{}
	var rs Rs

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.geekdo.com/api/geekmarket/products/%d", id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}
