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

	go monkey.KickOffSim()

	router.Run() // port 8080
}

func info(c *gin.Context) {
	c.JSON(http.StatusOK, monkey.FetchResults())
}
