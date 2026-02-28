package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blccming/private-positioning-server/docs"
	_ "github.com/blccming/private-positioning-server/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Session struct {
	Token      string    `json:"token"`
	Latitude   float32   `json:"latitude"`
	Longitude  float32   `json:"longitude"`
	TTL        int       `json:"ttl"`
	Timeout    int       `json:"timeout"`
	LastUpdate time.Time `json:"last_update"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid json"`
}

var sessions []Session

func initEndpoints() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/docs"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", getHealth)
	r.POST("/session/create", postSessionCreate)
	r.POST("/session/terminate", postSessionTerminate)
	return r
}

/* health  */

type HealthResponse struct {
	Status  string `json:"status" example:"OK"`
	Runtime string `json:"runtime" example:"38271.133967079s"`
}

// getHealth godoc
// @Summary      Health check
// @Description  Returns service health and runtime information
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func getHealth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "OK", "runtime": getRuntime()})
}

/* session create */

type SessionCreateRequest struct {
	TTL            int `json:"ttl" example:"3600"`
	SessionTimeout int `json:"session_timeout" example:"7200"`
}

type SessionCreateResponse struct {
	Token  string               `json:"token" example:"3A9N2O"`
	Params SessionCreateRequest `json:"parameters"`
}

// postSessionCreate godoc
// @Summary      Create a new session
// @Description  Creates a new session with the given TTL and session timeout
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Param        session  body      SessionCreateRequest  true  "Session create payload"
// @Success      200      {object}  SessionCreateResponse
// @Failure      400      {object}  ErrorResponse
// @Router       /sessions/create [post]
func postSessionCreate(c *gin.Context) {
	var newSession Session
	input := SessionCreateRequest{TTL: -1, SessionTimeout: -1}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if input.TTL <= 0 || input.SessionTimeout <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "ttl and session_timeout must be set and > 0"})
	}

	newSession.Token = tokenCreate()
	newSession.TTL = input.TTL
	newSession.Timeout = input.SessionTimeout

	// TODO: dont directly modify sessions[] -> no need to fix, will switch to redis later
	sessions = append(sessions, newSession)

	fmt.Println(newSession)
	c.JSON(http.StatusOK, gin.H{"token": newSession.Token})
}

/* terminate */
type SessionTerminateRequest struct {
	Token string `json:"token" example:"3A9N2O"`
}

// postSessionTerminate godoc
// @Summary      Terminate a session
// @Description  Terminates an active session by its token
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Param        session  body    SessionTerminateRequest  true  "Session terminate payload"
// @Success      200              {string}  string "successfully terminated session." TODO: fix example
// @Failure      400              {object}  ErrorResponse
// @Router       /sessions/terminate [post]
func postSessionTerminate(c *gin.Context) {
	var input struct {
		Token string `json:"token"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("%e", err))
		return
	}

	if !tokenExists(input.Token) {
		c.JSON(http.StatusBadRequest, "session token does not exist.")
	}

	index := sessionTokenToIndex(input.Token)
	if index == -1 {
		c.JSON(http.StatusBadRequest, "error while searching for Session struct.")
		return
	}
	sessions = append(sessions[:index], sessions[index+1:]...)

	c.JSON(http.StatusOK, "successfully terminated session.")
}
