package youtube

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"regexp"
)

func Handler(c *gin.Context) {
	url := c.Param("url")
	re := regexp.MustCompile(`youtube`)

	if !re.MatchString(url) {
		c.JSON(http.StatusOK, gin.H{"status": "not a youtube link"})
		return
	}

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
