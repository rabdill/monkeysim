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

	router.GET("/seats", seatedMonkeys)

	router.GET("/monkeys", allMonkeys)
	router.POST("/monkeys", addMonkey)
	router.PATCH("/monkeys/:id/stand", stand)
	router.PATCH("/monkeys/:id/sit", sit)

	go monkey.KickOffSim()

	router.Run() // port 8080
}

func seatedMonkeys(c *gin.Context) {
	c.JSON(http.StatusOK, monkey.FetchResults())
}

func allMonkeys(c *gin.Context) {
	c.JSON(http.StatusOK, monkey.FetchAll())
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

	err = monkey.Stand(id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.Status(202)
}

func sit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, err)
		return
	}

	err = monkey.Sit(id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.Status(202)
}
