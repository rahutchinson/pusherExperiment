package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Power string

const (
	On  Power = "on"
	Off       = "off"
)

type thing struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ThingType string `json:"type"`
	Status    status `json:"status"`
}

type status struct {
	TimeStamp   time.Time `json:"timestamp"`
	PowerStatus Power     `json:"power"`
}

var things = []thing{
	{ID: "338429", Name: "Ryan's", ThingType: "hovercraft", Status: status{TimeStamp: time.Now(), PowerStatus: On}},
	{ID: "558429", Name: "Mario's", ThingType: "spaceship", Status: status{TimeStamp: time.Now(), PowerStatus: Off}},
}

func getThings(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, things)
}

func addThing(c *gin.Context) {
	var newThing thing

	if err := c.BindJSON(&newThing); err != nil {
		return
	}

	things = append(things, newThing)
	c.IndentedJSON(http.StatusCreated, newThing)
}

func getThingById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range things {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "that thang don't exist"})
}

func main() {
	router := gin.Default()

	router.GET("/things", getThings)
	router.GET("/things/:id", getThingById)
	router.POST("/things", addThing)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
