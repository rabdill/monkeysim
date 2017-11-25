package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rabdill/monkeysim/monkey"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	go monkey.KickOffSim()

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"a": monkey.FetchResults(),
		})
	})

	router.Run() // port 8080
}
