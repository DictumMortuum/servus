package main

import (
	"html/template"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func loadTemplate() (*template.Template, error) {
	t := template.New("")

	for name, file := range Assets.Files {
		h, err := ioutil.ReadAll(file)

		if err != nil {
			return nil, err
		}

		t, err = t.New(name).Parse(string(h))

		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

func main() {
	r := gin.Default()

	t, err := loadTemplate()

	if err != nil {
		panic(err)
	}

	r.SetHTMLTemplate(t)
	r.Static("/assets", "./assets")

	r.GET("/startpage", func(c *gin.Context) {
		c.HTML(200, "/html/index.html", nil)
	})

	r.Run(":1234")
}
