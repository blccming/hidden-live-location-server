package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/blccming/hidden-live-location-server/db"
	"github.com/blccming/hidden-live-location-server/docs"
	_ "github.com/blccming/hidden-live-location-server/docs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	glide "github.com/valkey-io/valkey-glide/go/v2"
)

type ErrorResponse struct {
	Error string `json:"error" example:"invalid json"`
}

func InitEndpoints(useSwagger bool, dbClient *glide.Client) *gin.Engine {
	r := gin.Default()
	r.Use(PerClientRateLimit(5, 10)) // limit to 5 requests per second with burst of 10
	r.Use(GlobalRateLimit(100, 10))  // limit to 100 requests per second with burst of 10
	r.Use(MaxBodySize(1 << 10))      // max of 1 KB, biggest legitimate request should be ~ 500 bytes

	// swagger config
	if useSwagger {
		docs.SwaggerInfo.BasePath = "/"
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// our endpoints
	r.GET("/health", getHealth)
	r.POST("/session/create", func(c *gin.Context) { postSessionCreate(c, dbClient) })
	r.POST("/session/terminate", func(c *gin.Context) { postSessionTerminate(c, dbClient) })
	r.POST("/session/update", func(c *gin.Context) { postSessionUpdate(c, dbClient) })
	r.GET("/session/:token", func(c *gin.Context) { getSession(c, dbClient) })
	log.Info().Msg("Finalized Gin initialization.")

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
func postSessionCreate(c *gin.Context, g *glide.Client) {
	input := SessionCreateRequest{TTL: -1, SessionTimeout: -1}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		log.Error().Stack().Err(err).Msg("Failed to bind session create request")
		return
	}

	if input.TTL <= 0 || input.SessionTimeout <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ttl and session_timeout must be set and > 0"})
		log.Error().Msgf("Invalid TTL or session_timeout: ttl=%d, session_timeout=%d", input.TTL, input.SessionTimeout)
		return
	}

	token := tokenCreate(g)
	err := db.AddSession(g, token, input.SessionTimeout, input.TTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add session"})
		log.Error().Stack().Err(err).Msg("Failed to add session.")
		return
	}

	resp := SessionCreateResponse{
		Token:  token,
		Params: input,
	}
	c.JSON(http.StatusOK, resp)
	log.Info().Msgf("Session created: %s", token)
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
func postSessionTerminate(c *gin.Context, g *glide.Client) {
	var input struct {
		Token string `json:"token"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("%e", err)})
		log.Error().Stack().Err(err).Msgf("Failed to bind session terminate request: %s", input.Token)
		return
	}

	exists, err := db.SessionExists(g, input.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check session token existence"})
		log.Error().Stack().Err(err).Msgf("Failed to check session token existence: %s", input.Token)
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session token does not exist."})
		log.Error().Msgf("Session token does not exist: %s", input.Token)
		return
	}

	err = db.RemoveSession(g, input.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete session"})
		log.Error().Stack().Err(err).Msgf("Failed to delete session: %s", input.Token)
		return
	}

	resp := SessionTerminateResponse{
		Message: "successfully terminated session.",
		Token:   input.Token,
	}
	c.JSON(http.StatusOK, resp)
	log.Info().Msgf("Session terminated: %s", input.Token)
}

/*
 * session update
 */
type SessionUpdateRequest struct {
	Token     string  `json:"token" example:"3A9N2O"`
	Longitude float64 `json:"longitude" example:"49.026598"`
	Latitude  float64 `json:"latitude" example:"8.385259"`
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
func postSessionUpdate(c *gin.Context, g *glide.Client) {
	var input SessionUpdateRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("%e", err)})
		log.Error().Stack().Err(err).Msgf("Error while binding JSON")
		return
	}

	exists, err := db.SessionExists(g, input.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check session token existence"})
		log.Error().Stack().Err(err).Msgf("Failed to check session token existence: %s.", input.Token)
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session token does not exist."})
		log.Error().Msgf("Session token does not exist: %s", input.Token)
		return
	}

	err = db.SetLocation(g, input.Token, input.Longitude, input.Latitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update session location"})
		log.Error().Stack().Err(err).Msgf("Failed to update session location: %s", input.Token)
		return
	}

	resp := SessionUpdateResponse{
		Message: "successfully updated session.",
		Token:   input.Token,
	}
	c.JSON(http.StatusOK, resp)

	log.Info().Msgf("Session updated: %s", input.Token)
}

/*
 * get session
 */
type SessionGetResponse struct {
	Token      string    `json:"token" example:"3A9N2O"`
	Longitude  float64   `json:"longitude" example:"49.026598"`
	Latitude   float64   `json:"latitude" example:"8.385259"`
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
func getSession(c *gin.Context, g *glide.Client) {
	token := c.Param("token")

	exists, err := db.SessionExists(g, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check session token existence"})
		log.Error().Stack().Err(err).Msgf("Failed to check session token existence: %s.", token)
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session token does not exist."})
		log.Error().Msgf("Session token does not exist: %s.", token)
		return
	}

	locData, err := db.GetLocation(g, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get location data"})
		log.Error().Stack().Err(err).Msgf("Failed to get location data for session: %s", token)
		return
	}

	if locData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no location data available yet."})
		log.Error().Msgf("No location data available for session: %s", token)
		return
	}

	lonStr := locData["longitude"]
	latStr := locData["latitude"]
	lastChangedStr := locData["lastChanged"]

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse longitude"})
		log.Error().Stack().Err(err).Msgf("Failed to parse longitude for session: %s", token)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse latitude"})
		log.Error().Stack().Err(err).Msgf("Failed to parse latitude for session: %s", token)
		return
	}

	t, err := time.Parse(time.RFC3339, lastChangedStr) // or your format
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse lastChanged"})
		log.Error().Stack().Err(err).Msgf("Failed to parse lastChanged for session: %s", token)
		return
	}

	resp := SessionGetResponse{
		Token:      token,
		Longitude:  lon,
		Latitude:   lat,
		LastUpdate: t,
	}
	c.JSON(http.StatusOK, resp)

	log.Info().Msgf("Session retrieved: %s", token)
}
