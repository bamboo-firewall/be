package middleware

import (
	"github.com/bamboo-firewall/be/domain"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func Authorize(obj string, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		sub, existed := c.Get("x-user-id")
		if !existed {
			c.AbortWithStatusJSON(401, domain.ErrorResponse{Message: "User hasn't logged in yet"})
			return
		}

		// Load policy from Database
		err := enforcer.LoadPolicy()
		if err != nil {
			c.AbortWithStatusJSON(500, domain.ErrorResponse{Message: "Failed to load policy from DB"})
			return
		}

		// Casbin enforces policy
		ok, err := enforcer.Enforce(sub, obj, act)

		if err != nil {
			c.AbortWithStatusJSON(500, domain.ErrorResponse{Message: "Error occurred when authorizing user"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(403, domain.ErrorResponse{Message: "You are not authorized"})
			return
		}
		c.Next()
	}
}
