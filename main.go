package main

import (
	"github.com/gin-gonic/gin"
)

func startpageHandler(c *gin.Context) {
	data, _ := Asset("html/index.html")
	c.Writer.Write(data)
}

func main() {
	files := assetFS()
	r := gin.Default()
	r.StaticFS("/assets", files)
	r.GET("/startpage", startpageHandler)
	r.Run(":1234")
}
