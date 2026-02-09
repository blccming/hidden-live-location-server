package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type session struct {
	Token           string `json:"token"`
	TTL             int    `json:"ttl"`
	TerminationTime int    `json:"termination_time"`
}

func initEndpoints() *gin.Engine {
	router := gin.Default()
	router.GET("/health", getHealth)
	router.POST("/session/create", postSessionCreate)
	return router
}

func getHealth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "OK")
}

func postSessionCreate(c *gin.Context) {
	var newSession session
	if err := c.BindJSON(&newSession); err != nil {
		fmt.Println("%e", err)
		return
	}

	newSession.Token = "use-token-generation-here" // use 128-bit key?

	fmt.Println(newSession)
	c.IndentedJSON(http.StatusAccepted, newSession)
}
