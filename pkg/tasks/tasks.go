package tasks

import (
	"database/sql"
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
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type calendar struct {
	Id  int64  `db:"id"`
	Uri string `db:"principaluri"`
}

type props struct {
	Id    int64  `db:"objectid"`
	Name  string `db:"name"`
	Value string `db:"value"`
}

type Wish struct {
	Id         int64  `json:"Id" db:"id"`
	CalendarId int64  `json:"CalendarId" db:"calendar_id"`
	Owner      string `json:"Owner" db:"owner"`
	Status     string `json:"Status" db:"status"`
	Desc       string `json:"Desc" db:"description"`
	Title      string `json:"Title" db:"title"`
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

	err := db.Select(&rs, "select objectid, name, value from oc_calendarobjects_props where objectid = ?", id)
	if err != nil {
		return nil, err
	}

	for _, prop := range rs {
		switch prop.Name {
		case "SUMMARY":
			item.Id = prop.Id
			item.Title = prop.Value
		case "DESCRIPTION":
			item.Id = prop.Id
			item.Description = prop.Value
		case "STATUS":
			item.Id = prop.Id
			item.Status = prop.Value
		}
	}

	return &item, nil
}

func uriToUser(uri string) string {
	tmp := strings.Split(uri, "/")
	return tmp[2]
}

func exists(db *sqlx.DB, data Wish) (*sql.NullInt64, error) {
	var id sql.NullInt64

	stmt, err := db.PrepareNamed("select id from twishes where id = :id")
	if err != nil {
		return nil, err
	}

	err = stmt.Get(&id, data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &id, nil
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

func intInSlice(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func deleteRedundant(db *sqlx.DB, lists []List) error {
	rs := []int64{}
	local_rs := []int64{}

	err := db.Select(&rs, "select id from twishes")
	if err != nil {
		return err
	}

	for _, l := range lists {
		for _, item := range l.Items {
			local_rs = append(local_rs, item.Id)
		}
	}

	for _, item := range rs {
		if !intInSlice(item, local_rs) {
			_, err = db.Exec(`delete from twishes where id = $1`, item)
			if err != nil {
				return err
			}
		}
	}

	_, err = db.Exec(`delete from twishes where nextcloud_status = 'COMPLETED'`)
	if err != nil {
		return err
	}

	return nil
}

func syncList(db *sqlx.DB, list List) error {
	for _, item := range list.Items {
		payload := Wish{
			Id:         item.Id,
			CalendarId: list.Id,
			Owner:      list.User,
			Status:     item.Status,
			Desc:       item.Description,
			Title:      item.Title,
		}

		if payload.Status == "COMPLETED" {
			continue
		}

		id, err := exists(db, payload)
		if err != nil {
			return err
		}

		if id == nil {
			_, err = db.NamedExec(`
			insert into twishes (
				id,
				calendar_id,
				owner,
				status,
				description,
				title
			) values (
				:id,
				:calendar_id,
				:owner,
				:status,
				:description,
				:title
			)
		`, &payload)
			if err != nil {
				return err
			}
		} else {
			_, err = db.NamedExec(`
				update twishes set
					title = :title,
					description = :description,
					nextcloud_status = :status
				where
					id = :id and
					calendar_id = :calendar_id and
					owner = :owner
			`, &payload)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SyncTasks(c *gin.Context) {
	list := c.Param("list")

	db1, err := DB.DatabaseConnect("nextcloud")
	if err != nil {
		util.Error(c, err)
		return
	}
	defer db1.Close()

	db2, err := DB.DatabaseTypeConnect("postgres")
	if err != nil {
		util.Error(c, err)
		return
	}
	defer db2.Close()

	rs, err := getTasks(db1, list)
	if err != nil {
		util.Error(c, err)
		return
	}

	for _, l := range rs {
		err := syncList(db2, l)
		if err != nil {
			util.Error(c, err)
			return
		}
	}

	err = deleteRedundant(db2, rs)
	if err != nil {
		util.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, &rs)
}
