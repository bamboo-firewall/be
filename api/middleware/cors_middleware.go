package middleware

import (
	"github.com/bamboo-firewall/be/bootstrap"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware(env *bootstrap.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", env.CORSAllowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", env.CORSAllowOMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Continue to the next middleware
		c.Next()
	}
}
