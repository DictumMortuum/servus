package youtube

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Handler(c *gin.Context) {
	url := c.Param("url")

	f, err := os.OpenFile("/var/lib/servus/youtube.list", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	_, err = f.WriteString(url + "\n")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
