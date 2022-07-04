package bgg

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GeeklistItem struct {
	XMLName xml.Name `xml:"item" json:"-"`
	Id      int64    `xml:"objectid,attr" json:"id"`
	ItemId  int64    `xml:"id,attr" json:"item_id"`
	Name    string   `xml:"objectname,attr" json:"name"`
	Date    string   `xml:"editdate,attr" json:"date"`
	Body    string   `xml:"body" json:"body"`
}

type GeeklistRs struct {
	XMLName xml.Name       `xml:"geeklist" json:"-"`
	Items   []GeeklistItem `xml:"item" json:"items"`
}

func Geeklist(id int64) (*GeeklistRs, error) {

	// <item id="4474489" objecttype="thing" subtype="boardgame" objectid="42910" objectname="Peloponnes" username="giuseppelepew" postdate="Fri, 19 Feb 2016 20:45:33 +0000" editdate="Sat, 09 Apr 2022 07:11:30 +0000" thumbs="5" imageid="586432">
	// <body>
	// Εκδοση: Αγγλογερμανική 2η. Κατάσταση: Στην ζελατίνα Τιμή : 35 Στη ζελατίνα. Έχει βαθούλωμα στην μία γωνία(φωτό στο λινκ) https://www.boardgamegeek.com/geekmarket/product/824283
	// </body>
	// </item>
	// <item id="4504686" objecttype="thing" subtype="boardgame" objectid="153804" objectname="Dawgs of War" username="ikarosx" postdate="Sat, 05 Mar 2016 21:56:25 +0000" editdate="Wed, 17 Mar 2021 18:40:04 +0000" thumbs="0" imageid="1906078">
	// <body>
	// Εχει ανοιχτει η συσκευασια Τιμη 22Ε [url]https://boardgamegeek.com/geekmarket/product/899804[/url]
	// </body>
	// </item>
	// <item id="5135266" objecttype="thing" subtype="boardgame" objectid="95610" objectname="Warhammer: Invasion – The Eclipse of Hope" username="eraserhead" postdate="Tue, 10 Jan 2017 06:50:58 +0000" editdate="Wed, 17 Mar 2021 16:36:26 +0000" thumbs="0" imageid="955208">
	// <body>
	// [b]Κατάσταση[/b]: Tο κουτάκι δεν είναι ανοιγμένο. [b]Γλώσσα/'Εκδοση[/b]: Αγγλικά, FFG [b]Τιμή[/b]: EUR 10 + ό,τι μεταφορικά. https://www.boardgamegeek.com/geekmarket/product/1087535
	// </body>
	// </item>
	// <item id="5260031" objecttype="thing" subtype="boardgame" objectid="109215" objectname="Gunship: First Strike!" username="grunherz" postdate="Mon, 06 Mar 2017 21:19:02 +0000" editdate="Sat, 27 Mar 2021 16:29:20 +0000" thumbs="1" imageid="1179297">
	// <body>
	// Έκδοση: 2013 kickstarter αγγλική Κατάσταση: πολύ καλή Τιμή: 42€ Αστρομαχίες με θωρακισμένα σκάφη, αποστολές κρούσης κλπ. Laser, proton torpedoes, shields, fighter squadrons σε ημερήσια διάταξη. Σενάρια και μπόλικες εναλλακτικές με ένα κάρο πρόσθετο υλικό από τις επεκτάσεις [thing=121532][/thing], [thing=137559][/thing], [thing=121531][/thing][thing=121529][/thing], [thing=154757][/thing]. Με λίγα λόγια: Ameritrash experience on full afterburner https://www.boardgamegeek.com/geekmarket/product/1152248
	// </body>
	// </item>

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
