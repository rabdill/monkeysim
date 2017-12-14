package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rabdill/monkeysim/monkey"
)

func main() {
	router := gin.Default()

	router.StaticFile("/", "./static/index.html")
	router.Static("/static", "./static")

	router.GET("/monkeys", info)
	router.POST("/monkeys", addMonkey)

	router.PATCH("/monkeys/:id/stand", stand)

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

func stand(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, err)
		return
	}

	err = monkey.StandUp(id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.Status(202)
}
