package music

import (
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/fhs/gompd/mpd"
	"github.com/gin-gonic/gin"
	"net/http"
)

func stop(conn *mpd.Client) (*mpd.Attrs, error) {
	err := conn.Clear()
	if err != nil {
		return nil, err
	}

	err = conn.Repeat(false)
	if err != nil {
		return nil, err
	}

	err = conn.Stop()
	if err != nil {
		return nil, err
	}

	status, err := conn.Status()
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func playlist(conn *mpd.Client, list string) (*mpd.Attrs, error) {
	err := conn.Clear()
	if err != nil {
		return nil, err
	}

	err = conn.PlaylistLoad(list, -1, -1)
	if err != nil {
		return nil, err
	}

	err = conn.Repeat(true)
	if err != nil {
		return nil, err
	}

	err = conn.Play(-1)
	if err != nil {
		return nil, err
	}

	status, err := conn.Status()
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func Playlist(c *gin.Context) {
	list := c.Param("playlist")

	conn, err := mpd.Dial("tcp", config.App.Mpd.Server)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer conn.Close()

	status, err := conn.Status()
	if err != nil {
		util.Error(c, err)
		return
	}

	var retval *mpd.Attrs

	if status["state"] == "play" {
		err = conn.Next()
		if err != nil {
			util.Error(c, err)
			return
		}
	} else {
		retval, err = playlist(conn, list)
		if err != nil {
			util.Error(c, err)
			return
		}
	}

	util.Success(c, &retval)
}

func Stop(c *gin.Context) {
	conn, err := mpd.Dial("tcp", config.App.Mpd.Server)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer conn.Close()

	status, err := stop(conn)
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &status)
}

func Toggle(c *gin.Context) {
	conn, err := mpd.Dial("tcp", config.App.Mpd.Server)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer conn.Close()

	status, err := conn.Status()
	if err != nil {
		util.Error(c, err)
		return
	}

	if status["state"] == "play" {
		err = conn.Stop()
		if err != nil {
			util.Error(c, err)
			return
		}
	} else {
		err = conn.Play(-1)
		if err != nil {
			util.Error(c, err)
			return
		}
	}

	status, err = conn.Status()
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &status)
}

func Next(c *gin.Context) {
	conn, err := mpd.Dial("tcp", config.App.Mpd.Server)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer conn.Close()

	err = conn.Next()
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, nil)
}

func Previous(c *gin.Context) {
	conn, err := mpd.Dial("tcp", config.App.Mpd.Server)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer conn.Close()

	err = conn.Previous()
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, nil)
}

func Current(c *gin.Context) {
	conn, err := mpd.Dial("tcp", config.App.Mpd.Server)
	if err != nil {
		util.Error(c, err)
		return
	}
	defer conn.Close()

	attrs, err := conn.CurrentSong()
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &attrs)
}

func Radio(c *gin.Context) {
	state := map[string]interface{}{
		"calc": map[string]interface{}{
			"total_month": 0,
		},
	}

	c.HTML(http.StatusOK, "radio.html", state)
}
