package api

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

func InitEndpoints() *gin.Engine {
	r := gin.Default()
	// set middleware
	r.Use(RateLimit(5, 5))      // limit to 5 requests per second with burst of 5
	r.Use(MaxBodySize(1 << 10)) // max of 1 KB, biggest legitimate request should be ~ 500 bytes

	// swagger config
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// our endpoints
	r.GET("/health", getHealth)
	r.POST("/session/create", postSessionCreate)
	r.POST("/session/terminate", postSessionTerminate)
	r.POST("/session/update", postSessionUpdate)
	r.GET("/session/:token", getSession)

	return r
}

/*
 *
 * 	ENDPOINTS
 *
 */

/*
 * health check
 */

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

/*
 * session create
 */

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
// @Tags         session
// @Accept       json
// @Produce      json
// @Param        session  body      SessionCreateRequest  true  "Session create payload"
// @Success      200      {object}  SessionCreateResponse
// @Failure      400      {object}  ErrorResponse
// @Router       /session/create [post]
func postSessionCreate(c *gin.Context) {
	var newSession Session
	input := SessionCreateRequest{TTL: -1, SessionTimeout: -1}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if input.TTL <= 0 || input.SessionTimeout <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ttl and session_timeout must be set and > 0"})
		return
	}

	newSession.Token = tokenCreate()
	newSession.TTL = input.TTL
	newSession.Timeout = input.SessionTimeout

	// TODO: dont directly modify sessions[] -> no need to fix, will switch to redis later
	sessions = append(sessions, newSession)
	fmt.Println(newSession) // TODO: remove in production?

	resp := SessionCreateResponse{
		Token:  newSession.Token,
		Params: input,
	}
	c.JSON(http.StatusOK, resp)
}

/*
 * session terminate
 */
type SessionTerminateRequest struct {
	Token string `json:"token" example:"3A9N2O"`
}

type SessionTerminateResponse struct {
	Message string `json:"message" example:"successfully terminated session."`
	Token   string `json:"token" example:"3A9N2O"`
}

// postSessionTerminate godoc
// @Summary      Terminate a session
// @Description  Terminates an active session by its token
// @Tags         session
// @Accept       json
// @Produce      json
// @Param        session  body    SessionTerminateRequest  true  "Session terminate payload"
// @Success      200              {object}  SessionTerminateResponse "successfully terminated session."
// @Failure      400              {object}  ErrorResponse
// @Router       /session/terminate [post]
func postSessionTerminate(c *gin.Context) {
	var input struct {
		Token string `json:"token"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("%e", err)})
		return
	}

	if !tokenExists(input.Token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session token does not exist."})
		return
	}

	index := sessionTokenToIndex(input.Token)
	if index == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while searching for session struct."})
		return
	}
	sessions = append(sessions[:index], sessions[index+1:]...)

	resp := SessionTerminateResponse{
		Message: "successfully terminated session.",
		Token:   input.Token,
	}
	c.JSON(http.StatusOK, resp)
}

/*
 * session update
 */
type SessionUpdateRequest struct {
	Token     string  `json:"token" example:"3A9N2O"`
	Longitude float32 `json:"longitude" example:"49.026598"`
	Latitude  float32 `json:"latitude" example:"8.385259"`
}

type SessionUpdateResponse struct {
	Message string `json:"message" example:"successfully updated session."`
	Token   string `json:"token" example:"3A9N2O"`
}

// postSessionUpdate godoc
// @Summary      Update a session
// @Description  Updates the location information (longitude & latitude) of an active session identified by its token.
// @Tags         session
// @Accept       json
// @Produce      json
// @Param        session  body    SessionUpdateRequest  true  "Session update payload"
// @Success      200      {object} SessionUpdateResponse "Session updated."
// @Failure      400      {object} ErrorResponse
// @Router       /session/update [post]
func postSessionUpdate(c *gin.Context) {
	var input SessionUpdateRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("%e", err)})
		return
	}

	if !tokenExists(input.Token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session token does not exist."})
		return
	}

	index := sessionTokenToIndex(input.Token)
	if index == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while searching for session struct."})
		return
	}

	sessions[index].Longitude = input.Longitude
	sessions[index].Latitude = input.Latitude
	sessions[index].LastUpdate = time.Now()

	resp := SessionUpdateResponse{
		Message: "successfully updated session.",
		Token:   input.Token,
	}
	c.JSON(http.StatusOK, resp)
}

/*
 * get session
 */
type SessionGetResponse struct {
	Token      string    `json:"token" example:"3A9N2O"`
	Longitude  float32   `json:"longitude" example:"49.026598"`
	Latitude   float32   `json:"latitude" example:"8.385259"`
	LastUpdate time.Time `json:"last_update" example:"2026-03-01T13:23:45.206365244+01:00"`
}

// getSession godoc
// @Summary      Retrieve a session
// @Description  Fetches the current longitude, latitude and last‑update timestamp for a session identified by its token.
// @Tags         session
// @Accept       json
// @Produce      json
// @Param        token   path     string  true  "Session token"
// @Success      200     {object}  SessionGetResponse "Session retrieved."
// @Failure      400     {object}  ErrorResponse
// @Router       /session/{token} [get]
func getSession(c *gin.Context) {
	token := c.Param("token")
	fmt.Println(token)
	if !tokenExists(token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session token does not exist."})
		return
	}

	index := sessionTokenToIndex(token)
	if index == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while searching for session struct."})
		return
	}

	if sessions[index].LastUpdate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no location data available yet."})
		return
	} // TODO: after data expiration is implemented, make sure, this is still correct

	resp := SessionGetResponse{
		Token:      sessions[index].Token,
		Longitude:  sessions[index].Longitude,
		Latitude:   sessions[index].Latitude,
		LastUpdate: sessions[index].LastUpdate,
	}
	c.JSON(http.StatusOK, resp)
}
