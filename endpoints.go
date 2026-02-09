package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type session struct {
	Token           string    `json:"token"`
	TTL             int       `json:"ttl"`
	TerminationTime int       `json:"termination_time"`
	LastUpdate      time.Time `json:"last_update"`
}

var sessions []session

func initEndpoints() *gin.Engine {
	router := gin.Default()
	router.GET("/health", getHealth)
	router.POST("/session/create", postSessionCreate)
	return router
}

/* health  */

func getHealth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "OK", "runtime": getRuntime()})
}

/* session */

func postSessionCreate(c *gin.Context) {
	var newSession session

	// TODO: sanatize inputs
	if err := c.BindJSON(&newSession); err != nil {
		fmt.Println("%e", err)
		c.JSON(http.StatusBadRequest, fmt.Errorf("%e", err))
		return
	}

	newSession.Token = "use-token-generation-here" // use 128-bit key?

	// TODO: dont directly modify sessions[]
	sessions = append(sessions, newSession)

	fmt.Println(newSession)
	c.JSON(http.StatusOK, newSession)
}
