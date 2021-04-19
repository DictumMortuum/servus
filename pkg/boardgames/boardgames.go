package boardgames

import (
	"errors"
	DB "github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
)

func GetBoardgame(db *sqlx.DB, id int64) (interface{}, error) {
	rs, err := getBoardgame(db, id)
	if err != nil {
		return nil, err
	}

	rs.Atlas, err = AtlasSearch(*rs)
	if err != nil {
		return nil, err
	}

	for i, item := range rs.Atlas {
		tmp, err := BggRank(item)
		if err != nil {
			return nil, err
		}

		if len(tmp.Items) > 0 {
			log.Println(tmp.Items[0].Statistics.Ratings.Bayesaverage.Value)
			rs.Atlas[i].Rank = tmp.Items[0].Statistics.Ratings.Bayesaverage.Value
		}

		log.Println(item)

	}

	return rs, nil
}

func GetBoardgameList(db *sqlx.DB, args models.Args) (interface{}, int, error) {
	var rs []models.Boardgame

	sql, err := args.List(`
	select
		*
	from
		tboardgames
	`)
	if err != nil {
		return nil, -1, err
	}

	if len(args.Id) > 0 {
		query, ids, err := sqlx.In(sql.String(), args.Id)
		if err != nil {
			return nil, -1, err
		}

		err = db.Select(&rs, db.Rebind(query), ids...)
		if err != nil {
			return nil, -1, err
		}
	} else {
		err = db.Select(&rs, sql.String())
		if err != nil {
			return nil, -1, err
		}
	}

	return rs, len(rs), nil
}

func CreateBoardgame(db *sqlx.DB, data map[string]interface{}) (interface{}, error) {
	var player models.Boardgame

	if val, ok := data["name"]; ok {
		player.Name = val.(string)
	} else {
		return nil, errors.New("please provide a 'name' parameter")
	}

	sql := `
	insert into tboardgames (
		name
	) values (
		:name
	)`

	rs, err := db.NamedExec(sql, &player)
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}

	player.Id = id

	return player, nil
}

func UpdateBoardgame(db *sqlx.DB, id int64, data map[string]interface{}) (interface{}, error) {
	player, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	if val, ok := data["name"]; ok {
		player.Name = val.(string)
	}

	sql := `
	update
		tboardgames
	set
		name = :name
	where
		id = :id`

	_, err = db.NamedExec(sql, &player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func DeleteBoardgame(db *sqlx.DB, id int64) (interface{}, error) {
	player, err := getPlayer(db, id)
	if err != nil {
		return nil, err
	}

	sql := `
	delete from
		tboardgames
	where
		id = :id`

	_, err = db.NamedExec(sql, &player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func getBoardgame(db *sqlx.DB, id int64) (*models.Boardgame, error) {
	var retval models.Boardgame

	err := db.QueryRowx(`select * from tboardgames where id = ?`, id).StructScan(&retval)
	if err != nil {
		return nil, err
	}

	return &retval, nil
}

func updateAtlasId(ids, atlas_id string) error {
	database, err := DB.Conn()
	if err != nil {
		return err
	}
	defer database.Close()

	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return err
	}

	sql := `
	update
		tboardgames
	set
		atlas_id = ?
	where
		id = ?`

	_, err = database.Exec(sql, atlas_id, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateAtlasId(c *gin.Context) {
	arg1 := c.Params.ByName("id")
	arg2 := c.Params.ByName("atlas")
	log.Println(arg2)
	err := updateAtlasId(arg1, arg2)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Success(c, &arg1)
}

func Unmap(c *gin.Context) {
	arg1 := c.Params.ByName("id")

	err := unmap(arg1)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Success(c, &arg1)
}

func unmap(ids string) error {
	database, err := DB.Conn()
	if err != nil {
		return err
	}
	defer database.Close()

	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return err
	}

	sql := `
	update
		tboardgames
	set
		unmapped = 1
	where
		id = ?`

	_, err = database.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}
