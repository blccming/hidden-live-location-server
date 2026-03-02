package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

func RateLimit(perSec int, maxReq int) gin.HandlerFunc {
	var limiter = rate.NewLimiter(rate.Limit(perSec), maxReq)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			log.Warn().Msgf("Global rate limit exceeded: %d req/s, max %d req", perSec, maxReq)
			return
		}
		c.Next()
	}
}

func MaxBodySize(max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, max)
		log.Warn().Msgf("Max body size exceeded: %d bytes", max)
		c.Next()
	}
}
