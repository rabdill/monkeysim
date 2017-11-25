package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rabdill/monkeysim/monkey"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	monkey.KickOffSim()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"a": monkey.FetchResults(),
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
