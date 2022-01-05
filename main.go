package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pusher/pusher-http-go"
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
	{ID: "558429", Name: "Evgeny's", ThingType: "spaceship", Status: status{TimeStamp: time.Now(), PowerStatus: Off}},
}

func getThings(c *gin.Context) {
	pusherClient := pusher.Client{
		AppID:   "1326875",
		Key:     "d718acf39c8c6bfdddc6",
		Secret:  "d05fbeff194ba3e29a99",
		Cluster: "us2",
		Secure:  true,
	}
	pusherClient.Trigger("my-channel", "my-event", things)
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

func timePolling() {
	pusherClient := pusher.Client{
		AppID:   "1326875",
		Key:     "d718acf39c8c6bfdddc6",
		Secret:  "d05fbeff194ba3e29a99",
		Cluster: "us2",
		Secure:  true,
	}
	pusherClient.Trigger("my-channel", "my-event", things)
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/things", getThings)
	router.GET("/things/:id", getThingById)
	router.POST("/things", addThing)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
	for range time.Tick(time.Second * 10) {
		timePolling()
	}
}
