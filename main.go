package main

import (
	"github.com/DictumMortuum/servus/calendar"
	"github.com/DictumMortuum/servus/links"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func main() {
	mode := os.Getenv("GIN_MODE")

	if mode == "release" {
		gin.DisableConsoleColor()
		f, _ := os.Create("/var/log/servus.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}

	err := os.MkdirAll("/var/lib/servus", 0755)
	if err != nil {
		log.Println(err)
	}

	r := gin.Default()
	r.POST("/calendar", calendar.Handler)
	r.POST("/links", links.Handler)
	r.Run("127.0.0.1:1234")
}
