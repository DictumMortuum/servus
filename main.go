package main

import (
	"github.com/gin-gonic/gin"
)

func startpageHandler(c *gin.Context) {
	data, _ := Asset("html/index.html")
	c.Writer.Write(data)
}

func main() {
	r := gin.New()
	r.StaticFS("/assets", assetFS())
	r.GET("/startpage", startpageHandler)
	r.Run("127.0.0.1:1234")
}
