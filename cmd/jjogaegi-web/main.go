package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryanbrainard/jjogaegi/cmd"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	basePath := "cmd/jjogaegi-web"
	router.LoadHTMLGlob(basePath + "/templates/*.tmpl.html") // TODO: relative path
	router.Static("/static", basePath+"/static")
	indexTemplate := "index.tmpl.html"

	router.GET("/", func(c *gin.Context) {
		var m Model
		m.init()

		if err := c.Bind(&m); err != nil {
			log.Printf("at=get.bind error=%q", err)
			m.Error = fmt.Sprint(err)
			c.HTML(http.StatusBadRequest, indexTemplate, m)
			return
		}

		c.HTML(http.StatusOK, indexTemplate, m)
	})

	router.POST("/", func(c *gin.Context) {
		var m Model
		m.init()

		if err := c.Bind(&m); err != nil {
			log.Printf("at=post.bind error=%q", err)
			m.Error = err.Error()
			c.HTML(http.StatusBadRequest, indexTemplate, m)
			return
		}

		output := &bytes.Buffer{}

		err := pkg.Run(
			strings.NewReader(m.Input),
			output,
			cmd.ParseOptParser(m.Parser),
			cmd.ParseOptFormatter(m.Formatter),
			map[string]string{},
		)

		if err != nil {
			log.Printf("at=post.run error=%q", err)
			m.Error = err.Error()
			c.HTML(http.StatusBadRequest, indexTemplate, m)
			return
		}

		m.Output = output.String()
		c.HTML(http.StatusOK, indexTemplate, m)
	})

	router.Run(":" + port)
}

// TODO: rename
type Model struct {
	Input        string `form:"input"`
	Output       string
	Error        string
	Parser       string `form:"parser"`
	Formatter    string `form:"formatter"`
	Capabilities cmd.Capabilities
}

func (m *Model) init() {
	m.Capabilities = cmd.AppCapabilities
}
