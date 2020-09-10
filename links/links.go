package links

import (
	"github.com/DictumMortuum/servus/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler(c *gin.Context) {
	url := c.PostForm("url")

	db, err := db.Conn()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into tlink (url) values (?)")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
