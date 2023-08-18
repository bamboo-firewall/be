package middleware

import (
	"github.com/bamboo-firewall/be/domain"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authorize(obj string, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		sub, existed := c.Get("x-user-id")
		if !existed {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "The request unauthorized"})
			return
		}

		// Load policy from Database
		err := enforcer.LoadPolicy()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to load policy"})
			return
		}

		// Casbin enforces policy
		ok, err := enforcer.Enforce(sub, obj, act)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Error occurred when authorizing user"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, domain.ErrorResponse{Message: "You are not authorized"})
			return
		}
		c.Next()
	}
}
