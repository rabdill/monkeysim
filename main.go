package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rabdill/monkeysim/monkey"
)

func main() {
	router := gin.Default()

	router.StaticFile("/", "./static/index.html")
	router.Static("/static", "./static")

	router.GET("/info", info)
	router.GET("/add_monkey", addMonkey)

	go monkey.KickOffSim()

	router.Run() // port 8080
}

func info(c *gin.Context) {
	c.JSON(http.StatusOK, monkey.FetchResults())
}

func addMonkey(c *gin.Context) {
	monkey, err := monkey.AddMonkey()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(201, monkey)
}
