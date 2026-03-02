package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

var clientLimiters = make(map[string]*rate.Limiter)

func PerClientRateLimit(perSec int, maxReq int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if _, exists := clientLimiters[ip]; !exists {
			clientLimiters[ip] = rate.NewLimiter(rate.Limit(perSec), maxReq)
		}

		if !clientLimiters[ip].Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			log.Warn().Msgf("Client rate limit exceeded: %s, %d req/s, max %d req", ip, perSec, maxReq)
			return
		}
		c.Next()
	}
}

func GlobalRateLimit(perSec int, maxReq int) gin.HandlerFunc {
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
		c.Next()
	}
}
