package boardgames

import (
	"encoding/xml"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"io/ioutil"
	"net/http"
)

type Bayes struct {
	XMLName xml.Name `xml:"bayesaverage"`
	Value   string   `xml:"value,attr"`
}

type Ratings struct {
	XMLName      xml.Name `xml:"ratings"`
	Bayesaverage Bayes    `xml:"bayesaverage"`
}

type Stats struct {
	XMLName xml.Name `xml:"statistics"`
	Ratings Ratings  `xml:"ratings"`
}

type Item struct {
	XMLName    xml.Name `xml:"item"`
	Statistics Stats    `xml:"statistics"`
}

type rs struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

func BggRank(game models.AtlasResult) (*rs, error) {

	// 	<statistics page="1">
	// 	<ratings >
	// 		<usersrated value="61737" />
	// 		<average value="8.10714" />
	// 		<bayesaverage value="7.97995" />
	// 		<ranks>
	// 			<rank type="subtype" id="1" name="boardgame" friendlyname="Board Game Rank" value="17" bayesaverage="7.97995" />
	// 			<rank type="family" id="5497" name="strategygames" friendlyname="Strategy Game Rank" value="18" bayesaverage="7.94971" />
	// 		</ranks>
	// 		<stddev value="1.18003" />
	// 		<median value="0" />
	// 		<owned value="96658" />
	// 		<trading value="625" />
	// 		<wanting value="955" />
	// 		<wishing value="8431" />
	// 		<numcomments value="8512" />
	// 		<numweights value="1921" />
	// 		<averageweight value="2.2233" />
	// 	</ratings>
	// </statistics>
	// </item>
	// </items>

	link := fmt.Sprintf("https://www.boardgamegeek.com/xmlapi2/thing?id=%d&stats=1", game.BggId)
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

	tmp := rs{}

	err = xml.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}

	return &tmp, nil
}
