package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type session struct {
	Token      string    `json:"token"`
	Latitude   float32   `json:"latitude"`
	Longitude  float32   `json:"longitude"`
	TTL        int       `json:"ttl"`
	Timeout    int       `json:"timeout"`
	LastUpdate time.Time `json:"last_update"`
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
	c.IndentedJSON(http.StatusOK, gin.H{"statunewSessions": "OK", "runtime": getRuntime()})
}

/* session */

func postSessionCreate(c *gin.Context) {
	var newSession session

	var input struct {
		TTL            int `json:"ttl"`
		SessionTimeout int `json:"session_timeout"`
	}

	// TODO: sanatize inputs
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("%e", err))
		return
	}

	newSession.Token = tokenCreate() // use 128-bit key
	newSession.TTL = input.TTL
	newSession.Timeout = input.SessionTimeout

	// TODO: dont directly modify sessions[] -> no need to fix, will switch to redis later
	sessions = append(sessions, newSession)

	fmt.Println(newSession)
	c.JSON(http.StatusOK, newSession)
}

func postSessionTerminate(c *gin.Context) {
	var input struct {
		Token string `json:"token"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("%e", err))
		return
	}

	if !tokenExists(input.Token) {
		c.JSON(http.StatusBadRequest, "Session token does not exist.")
	}

	index := sessionTokenToIndex(input.Token)
	if index == -1 {
		c.JSON(http.StatusBadRequest, "Error while searching for Session struct.")
		return
	}
	sessions = append(sessions[:index], sessions[index+1:]...)

	c.JSON(http.StatusOK, "Successfully terminated session.")
}
