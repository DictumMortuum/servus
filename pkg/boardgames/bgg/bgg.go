package bgg

import (
	"encoding/xml"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
)

type Rank struct {
	XMLName xml.Name `xml:"rank" json:"-"`
	Type    string   `xml:"name,attr" json:"type"`
	Name    string   `xml:"friendlyname,attr" json:"name"`
	Value   string   `xml:"value,attr" json:"value"`
}

type Bayes struct {
	XMLName xml.Name `xml:"bayesaverage" json:"-"`
	Value   string   `xml:"value,attr" json:"value"`
}

type Average struct {
	XMLName xml.Name `xml:"average" json:"-"`
	Value   string   `xml:"value,attr" json:"value"`
}

type Ranks struct {
	XMLName xml.Name `xml:"ranks" json:"-"`
	Ranks   []Rank   `xml:"rank" json:"ranks"`
}

type Ratings struct {
	XMLName      xml.Name `xml:"ratings" json:"-"`
	Bayesaverage Bayes    `xml:"bayesaverage"`
	Average      Average  `xml:"average"`
	Ranks        Ranks    `xml:"ranks" json:"ranks"`
}

type BggStats struct {
	XMLName xml.Name `xml:"statistics" json:"-"`
	Ratings Ratings  `xml:"ratings" json:"ratings"`
}

type BggName struct {
	XMLName xml.Name `xml:"name" json:"-"`
	Value   string   `xml:"value,attr" json:"value"`
}

type Item struct {
	XMLName    xml.Name `xml:"item" json:"-"`
	Statistics BggStats `xml:"statistics"`
	Name       BggName  `xml:"name"`
	Thumbnail  string   `xml:"thumbnail"`
	Image      string   `xml:"image"`
}

type BggThing struct {
	XMLName xml.Name `xml:"items" json:"-"`
	Items   []Item   `xml:"item"`
}

type rs struct {
	XMLName xml.Name `xml:"items" json:"-"`
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

// <?xml version="1.0" encoding="utf-8"?>
// <items total="1" termsofuse="https://boardgamegeek.com/xmlapi/termsofuse">
// 	<item type="boardgame" id="237182">
//   	<name type="primary" value="Root"/>
//     <yearpublished value="2018" />
//   </item>
// </items>

type SearchItem struct {
	XMLName xml.Name `xml:"item" json:"-"`
	Id      int64    `xml:"id,attr" json:"id"`
	Name    struct {
		XMLName xml.Name `xml:"name" json:"-"`
		Value   string   `xml:"value,attr" json:"value"`
	} `json:"name"`
	Published struct {
		XMLName xml.Name `xml:"yearpublished" json:"-"`
		Value   string   `xml:"value,attr" json:"value"`
	} `json:"published"`
}

type SearchRs struct {
	XMLName xml.Name     `xml:"items" json:"-"`
	Items   []SearchItem `xml:"item" json:"items"`
}

func filterItems(col []SearchItem) []SearchItem {
	rs := []SearchItem{}
	tmp := map[int64]SearchItem{}

	for _, item := range col {
		tmp[item.Id] = item
	}

	for key := range tmp {
		rs = append(rs, tmp[key])
	}

	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Id < rs[j].Id
	})

	return rs
}

func Search(name string) ([]SearchItem, error) {
	link := fmt.Sprintf("https://www.boardgamegeek.com/xmlapi2/search?query=%s&type=boardgameaccessory,boardgame,boardgameexpansion", url.QueryEscape(name))

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

	rs := SearchRs{
		Items: []SearchItem{},
	}
	err = xml.Unmarshal(body, &rs)
	if err != nil {
		return nil, err
	}

	return filterItems(rs.Items), nil
}

func BggSearch(c *gin.Context) {
	type params struct {
		Term string `json:"term"`
	}

	var payload params
	c.BindJSON(&payload)

	link := fmt.Sprintf("https://www.boardgamegeek.com/xmlapi2/search?query=%s&type=boardgame", url.QueryEscape(payload.Term))
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		util.Error(c, err)
		return
	}

	conn := &http.Client{}

	resp, err := conn.Do(req)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error(c, err)
		return
	}

	rs := SearchRs{}

	err = xml.Unmarshal(body, &rs)
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}

// <?xml version="1.0" encoding="utf-8"?>
// <items termsofuse="https://boardgamegeek.com/xmlapi/termsofuse">
// 	<item type="boardgame" id="3076">
// 		<thumbnail>https://cf.geekdo-images.com/QFiIRd2kimaMqTyWsX0aUg__thumb/img/5fLo89ChZH6Wzukk36bhZ-EpBS0=/fit-in/200x150/filters:strip_icc()/pic158548.jpg</thumbnail>
// 		<image>https://cf.geekdo-images.com/QFiIRd2kimaMqTyWsX0aUg__original/img/DOgIp57F7tKZvxeITGAd3e_Q9as=/0x0/filters:format(jpeg)/pic158548.jpg</image>
// 		<name type="primary" sortindex="1" value="Puerto Rico" />
// 		<name type="alternate" sortindex="1" value="Πουέρτο Ρίκο" />
// 		<name type="alternate" sortindex="1" value="波多黎各" />
// 		<name type="alternate" sortindex="1" value="푸에르토 리코" />
// 		<description>In Puerto Rico, players assume the roles of colonial governors on the island of Puerto Rico. The aim of the game is to amass victory points by shipping goods to Europe or by constructing buildings.&amp;#10;&amp;#10;Each player uses a separate small board with spaces for city buildings, plantations, and resources. Shared between the players are three ships, a trading house, and a supply of resources and doubloons.&amp;#10;&amp;#10;The resource cycle of the game is that players grow crops which they exchange for points or doubloons. Doubloons can then be used to buy buildings, which allow players to produce more crops or give them other abilities. Buildings and plantations do not work unless they are manned by colonists.&amp;#10;&amp;#10;During each round, players take turns selecting a role card from those on the table (such as &amp;quot;Trader&amp;quot; or &amp;quot;Builder&amp;quot;). When a role is chosen, every player gets to take the action appropriate to that role. The player that selected the role also receives a small privilege for doing so - for example, choosing the &amp;quot;Builder&amp;quot; role allows all players to construct a building, but the player who chose the role may do so at a discount on that turn. Unused roles gain a doubloon bonus at the end of each turn, so the next player who chooses that role gets to keep any doubloon bonus associated with it. This encourages players to make use of all the roles throughout a typical course of a game.&amp;#10;&amp;#10;Puerto Rico uses a variable phase order mechanism in which a &amp;quot;governor&amp;quot; token is passed clockwise to the next player at the conclusion of a turn. The player with the token begins the round by choosing a role and taking the first action.&amp;#10;&amp;#10;Players earn victory points for owning buildings, for shipping goods, and for manned &amp;quot;large buildings.&amp;quot; Each player's accumulated shipping chips are kept face down and come in denominations of one or five. This prevents other players from being able to determine the exact score of another player. Goods and doubloons are placed in clear view of other players and the totals of each can always be requested by a player. As the game enters its later stages, the unknown quantity of shipping tokens and its denominations require players to consider their options before choosing a role that can end the game.&amp;#10;&amp;#10;In 2011 and mostly afterwards, Puerto Rico was published to include both Puerto Rico: Expansion I &amp;ndash; New Buildings and Puerto Rico: Expansion II &amp;ndash; The Nobles. These versions are included in the other game entry Puerto Rico (with two expansions), not this regular game entry for Puerto Rico. Some editions of Puerto Rico list the player count as 2-5 instead of 3-5, and they include variant rules for games with only two players.&amp;#10;&amp;#10;</description>
// 		<yearpublished value="2002" />
// 		<minplayers value="3" />
// 		<maxplayers value="5" />
// 		<poll name="suggested_numplayers" title="User Suggested Number of Players" totalvotes="1464">
// 			<results numplayers="1">
// 				<result value="Best" numvotes="5" />
// 				<result value="Recommended" numvotes="11" />
// 				<result value="Not Recommended" numvotes="792" />
// 			</results>
// 			<results numplayers="2">
// 				<result value="Best" numvotes="78" />
// 				<result value="Recommended" numvotes="462" />
// 				<result value="Not Recommended" numvotes="500" />
// 			</results>
// 			<results numplayers="3">
// 				<result value="Best" numvotes="338" />
// 				<result value="Recommended" numvotes="865" />
// 				<result value="Not Recommended" numvotes="59" />
// 			</results>
// 			<results numplayers="4">
// 				<result value="Best" numvotes="968" />
// 				<result value="Recommended" numvotes="355" />
// 				<result value="Not Recommended" numvotes="5" />
// 			</results>
// 			<results numplayers="5">
// 				<result value="Best" numvotes="365" />
// 				<result value="Recommended" numvotes="689" />
// 				<result value="Not Recommended" numvotes="143" />
// 			</results>
// 			<results numplayers="5+">
// 				<result value="Best" numvotes="6" />
// 				<result value="Recommended" numvotes="23" />
// 				<result value="Not Recommended" numvotes="621" />
// 			</results>
// 		</poll>
// 		<playingtime value="150" />
// 		<minplaytime value="90" />
// 		<maxplaytime value="150" />
// 		<minage value="12" />
// 		<poll name="suggested_playerage" title="User Suggested Player Age" totalvotes="308">
// 			<results>
// 				<result value="2" numvotes="1" />
// 				<result value="3" numvotes="0" />
// 				<result value="4" numvotes="0" />
// 				<result value="5" numvotes="1" />
// 				<result value="6" numvotes="2" />
// 				<result value="8" numvotes="14" />
// 				<result value="10" numvotes="71" />
// 				<result value="12" numvotes="127" />
// 				<result value="14" numvotes="72" />
// 				<result value="16" numvotes="15" />
// 				<result value="18" numvotes="4" />
// 				<result value="21 and up" numvotes="1" />
// 			</results>
// 		</poll>
// 		<poll name="language_dependence" title="Language Dependence" totalvotes="320">
// 			<results>
// 				<result level="6" value="No necessary in-game text" numvotes="9" />
// 				<result level="7" value="Some necessary text - easily memorized or small crib sheet" numvotes="112" />
// 				<result level="8" value="Moderate in-game text - needs crib sheet or paste ups" numvotes="180" />
// 				<result level="9" value="Extensive use of text - massive conversion needed to be playable" numvotes="18" />
// 				<result level="10" value="Unplayable in another language" numvotes="1" />
// 			</results>
// 		</poll>
// 		<link type="boardgamecategory" id="1029" value="City Building" />
// 		<link type="boardgamecategory" id="1021" value="Economic" />
// 		<link type="boardgamecategory" id="1013" value="Farming" />
// 		<link type="boardgamemechanic" id="2838" value="Action Drafting" />
// 		<link type="boardgamemechanic" id="2875" value="End Game Bonuses" />
// 		<link type="boardgamemechanic" id="2843" value="Follow" />
// 		<link type="boardgamemechanic" id="2987" value="Hidden Victory Points" />
// 		<link type="boardgamemechanic" id="2914" value="Increase Value of Unchosen Resources" />
// 		<link type="boardgamemechanic" id="2828" value="Turn Order: Progressive" />
// 		<link type="boardgamemechanic" id="2079" value="Variable Phase Order" />
// 		<link type="boardgamefamily" id="10656" value="Country: Puerto Rico" />
// 		<link type="boardgamefamily" id="70360" value="Digital Implementations: Board Game Arena" />
// 		<link type="boardgamefamily" id="9958" value="Game: Puerto Rico" />
// 		<link type="boardgamefamily" id="27646" value="Mechanism: Tableau Building" />
// 		<link type="boardgamefamily" id="58" value="Series: Alea Big Box" />
// 		<link type="boardgamefamily" id="21622" value="Theme: Colonial" />
// 		<link type="boardgamefamily" id="18319" value="Theme: Tropical" />
// 		<link type="boardgameexpansion" id="286086" value="Deutscher Spielepreis Classic Goodie Box" />
// 		<link type="boardgameexpansion" id="69178" value="Haiti (fan expansion for Puerto Rico)" />
// 		<link type="boardgameexpansion" id="292912" value="Puerto Rico: Buccaneer" />
// 		<link type="boardgameexpansion" id="12382" value="Puerto Rico: Expansion I – New Buildings" />
// 		<link type="boardgameexpansion" id="159761" value="Puerto Rico: Expansion II – The Nobles" />
// 		<link type="boardgameexpansion" id="270125" value="Puerto Rico: Expansions 1&amp;2 – The New Buildings &amp; The Nobles" />
// 		<link type="boardgameexpansion" id="291423" value="Puerto Rico: Festival en Puerto Rico" />
// 		<link type="boardgameexpansion" id="40688" value="Treasure Chest" />
// 		<link type="boardgamecompilation" id="318985" value="Puerto Rico" />
// 		<link type="boardgamecompilation" id="108687" value="Puerto Rico (with two expansions)" />
// 		<link type="boardgamecompilation" id="165332" value="Puerto Rico + Expansion I" />
// 		<link type="boardgameimplementation" id="8217" value="San Juan" />
// 		<link type="boardgameimplementation" id="166669" value="San Juan (Second Edition)" />
// 		<link type="boardgamedesigner" id="117" value="Andreas Seyfarth" />
// 		<link type="boardgameartist" id="4959" value="Harald Lieske" />
// 		<link type="boardgameartist" id="11883" value="Franz Vohwinkel" />
// 		<link type="boardgamepublisher" id="9" value="alea" />
// 		<link type="boardgamepublisher" id="34" value="Ravensburger" />
// 		<link type="boardgamepublisher" id="4304" value="Albi" />
// 		<link type="boardgamepublisher" id="2366" value="Devir" />
// 		<link type="boardgamepublisher" id="5657" value="Filosofia Éditions" />
// 		<link type="boardgamepublisher" id="6214" value="Kaissa Chess &amp; Games" />
// 		<link type="boardgamepublisher" id="5812" value="Lacerta" />
// 		<link type="boardgamepublisher" id="3218" value="Lautapelit.fi" />
// 		<link type="boardgamepublisher" id="7992" value="MINDOK" />
// 		<link type="boardgamepublisher" id="3" value="Rio Grande Games" />
// 		<link type="boardgamepublisher" id="3888" value="Stratelibri" />
// 		<link type="boardgamepublisher" id="9234" value="Swan Panasia Co., Ltd." />
// 		<link type="boardgamepublisher" id="42" value="Tilsit" />
// 		<link type="boardgamepublisher" id="46567" value="VAKKO" />
// 		<statistics page="1">
// 			<ratings >
// 				<usersrated value="63640" />
// 				<average value="7.98406" />
// 				<bayesaverage value="7.85737" />
// 				<ranks>
// 					<rank type="subtype" id="1" name="boardgame" friendlyname="Board Game Rank" value="28" bayesaverage="7.85737" />
// 					<rank type="family" id="5497" name="strategygames" friendlyname="Strategy Game Rank" value="27" bayesaverage="7.84343" />
// 				</ranks>
// 				<stddev value="1.38235" />
// 				<median value="0" />
// 				<owned value="74012" />
// 				<trading value="1330" />
// 				<wanting value="1003" />
// 				<wishing value="9241" />
// 				<numcomments value="12032" />
// 				<numweights value="6007" />
// 				<averageweight value="3.2783" />
// 			</ratings>
// 		</statistics>
// 	</item>
// </items>

type GetItem struct {
	XMLName   xml.Name `xml:"item" json:"-"`
	Id        int64    `xml:"id,attr" json:"id"`
	Thumbnail string   `xml:"thumbnail" json:"thumb"`
	Image     string   `xml:"image" json:"image"`
	Name      struct {
		XMLName xml.Name `xml:"name" json:"-"`
		Value   string   `xml:"value,attr" json:"value"`
	} `json:"-"`
	Statistics struct {
		XMLName xml.Name `xml:"statistics" json:"-"`
		Ratings struct {
			XMLName      xml.Name `xml:"ratings" json:"-"`
			Ranks        Ranks    `xml:"ranks" json:"ranks"`
			Bayesaverage Bayes    `xml:"bayesaverage" json:"bayes"`
			Average      Average  `xml:"average" json:"average"`
		} `json:"ratings"`
	} `json:"published"`
}

type GetRs struct {
	XMLName xml.Name  `xml:"items" json:"-"`
	Items   []GetItem `xml:"item" json:"items"`
}

func BggGet(c *gin.Context) {
	type params struct {
		Id int64 `json:"id"`
	}

	var payload params
	c.BindJSON(&payload)

	link := fmt.Sprintf("https://www.boardgamegeek.com/xmlapi2/thing?id=%d&stats=1", payload.Id)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		util.Error(c, err)
		return
	}

	conn := &http.Client{}

	resp, err := conn.Do(req)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error(c, err)
		return
	}

	rs := GetRs{}

	err = xml.Unmarshal(body, &rs)
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &rs)
}
