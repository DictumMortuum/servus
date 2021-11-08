package tasks

import (
	DB "github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strings"
)

type List struct {
	Id    int64  `json:"id"`
	User  string `json:"user"`
	Items []Item `json:"items"`
}

type Item struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type calendar struct {
	Id  int64  `db:"id"`
	Uri string `db:"principaluri"`
}

type props struct {
	Name  string `db:"name"`
	Value string `db:"value"`
}

func getTasks(db *sqlx.DB, name string) ([]List, error) {
	rs := []calendar{}
	retval := []List{}

	err := db.Select(&rs, "select id, principaluri from oc_calendars where uri = ?", name)
	if err != nil {
		return nil, err
	}

	for _, calendar := range rs {
		list := List{}
		list.Id = calendar.Id
		list.User = uriToUser(calendar.Uri)
		list.Items, err = getProps(db, calendar.Id)
		if err != nil {
			return nil, err
		}

		retval = append(retval, list)
	}

	return retval, nil
}

func getProps(db *sqlx.DB, calendarid int64) ([]Item, error) {
	rs := []int64{}
	retval := []Item{}

	err := db.Select(&rs, "select distinct objectid from oc_calendarobjects_props where calendarid = ?", calendarid)
	if err != nil {
		return nil, err
	}

	for _, id := range rs {
		item, err := propsToItem(db, id)
		if err != nil {
			return nil, err
		}

		retval = append(retval, *item)
	}

	return retval, nil
}

func propsToItem(db *sqlx.DB, id int64) (*Item, error) {
	item := Item{}
	rs := []props{}

	err := db.Select(&rs, "select name, value from oc_calendarobjects_props where objectid = ?", id)
	if err != nil {
		return nil, err
	}

	for _, prop := range rs {
		switch prop.Name {
		case "SUMMARY":
			item.Title = prop.Value
		case "DESCRIPTION":
			item.Description = prop.Value
		case "STATUS":
			item.Status = prop.Value
		}
	}

	return &item, nil
}

func uriToUser(uri string) string {
	tmp := strings.Split(uri, "/")
	return tmp[2]
}

func GetTasks(c *gin.Context) {
	list := c.Param("list")

	db, err := DB.DatabaseConnect("nextcloud")
	if err != nil {
		util.Error(c, err)
		return
	}
	defer db.Close()

	rs, err := getTasks(db, list)
	if err != nil {
		util.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, &rs)
}
