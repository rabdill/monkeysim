package main

import (
	"fmt"
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
	router.POST("/seats", addSeat)
	router.PATCH("/seats/:id/stand", stand) // You tell a seat to stand, but a monkey to sit

	router.GET("/monkeys", allMonkeys)
	router.POST("/monkeys", addMonkey)

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

func addSeat(c *gin.Context) {
	var json monkey.AddSeatInput
	if err := c.ShouldBindJSON(&json); err != nil {
		fmt.Printf("\nERROR IN HERE 1 : %v", err)
		c.JSON(400, err.Error())
		return
	}
	seat, err := monkey.AddSeat(json)
	if err != nil {
		fmt.Printf("\nERROR IN HERE 2 : %v", err)
		c.JSON(500, err.Error())
		return
	}
	c.JSON(201, seat)
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
